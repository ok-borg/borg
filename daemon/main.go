package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"time"

	log "github.com/cihub/seelog"
	"github.com/crufter/borg/daemon/auth"
	"github.com/crufter/borg/types"
	"github.com/crufter/slugify"
	"github.com/joeguo/sitemap"
	httpr "github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"golang.org/x/oauth2"
	"gopkg.in/olivere/elastic.v2"
)

const (
	githubTokenURL = "https://github.com/login/oauth/access_token"
)

var (
	esAddr             = flag.String("esaddr", "127.0.0.1:9200", "Elastic Search address")
	githubClientId     = flag.String("github-client-id", "", "Github oauth client id")
	githubClientSecret = flag.String("github-client-secret", "", "Github client secret")
	sm                 = flag.String("sitemap", "", "Sitemap location. Leave empty if you don't want a sitemap to be generated")
)

var (
	client                 *elastic.Client
	aut                    *auth.Auth
	accessControl          map[string]UserAccess
	mtx                    = &sync.Mutex{}
	lastAccessControlReset = time.Now()
)

type AccessKinds int

// FIXME(jeremy): should be in config
// maximum access for write and updates
const (
	maxCreate = 100
	maxUpdate = 50
)

// acces kings
const (
	Create AccessKinds = iota
	Update
)

type Logger struct{}

type UserAccess struct {
	Update int
	Create int
}

func (l Logger) Printf(str string, i ...interface{}) {
	fmt.Println(fmt.Sprintf(str, i...))
}

func init() {
	flag.Parse()
	cl, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(fmt.Sprintf("http://%v", *esAddr)), elastic.SetTraceLog(Logger{}))
	if err != nil {
		panic(err)
	}
	client = cl
	accessControl = map[string]UserAccess{}
}

func main() {
	oauthCfg := &oauth2.Config{
		ClientID:     *githubClientId,
		ClientSecret: *githubClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: githubTokenURL,
		},
		Scopes: []string{"read:org"},
	}
	aut = auth.NewAuth(oauthCfg, client)
	r := httpr.New()
	if len(*sm) > 0 {
		go sitemapLoop()
	}

	r.GET("/v1/redirect/github/authorize", redirectGithubAuthorize)
	r.GET("/v1/query", query)
	r.POST("/v1/auth/github", githubAuth)

	// authenticated endpoints
	r.GET("/v1/user", ifAuth(getUser))

	// snippets
	r.GET("/v1/p/:id", getSnippet)
	r.GET("/v1/latest", getLatestSnippets)
	r.POST("/v1/p", ifAuth(controlAccess(createSnippet, Create)))
	r.DELETE("/v1/p/:id", ifAuth(deleteSnippet))
	r.PUT("/v1/p", ifAuth(controlAccess(updateSnippet, Update)))

	handler := cors.New(cors.Options{AllowedHeaders: []string{"*"}, AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"}}).Handler(r)
	log.Info("Starting http server")
	log.Critical(http.ListenAndServe(fmt.Sprintf(":%v", 9992), handler))
}

func writeJsonResponse(w http.ResponseWriter, status int, body interface{}) {
	rawBody, _ := json.Marshal(body)
	writeResponse(w, status, string(rawBody))
}

func writeResponse(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(body)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `%v`, body)
}

func updateTimer() {
	mtx.Lock()
	// if last reset was 24 or more ago
	if time.Since(lastAccessControlReset) >= (time.Hour * 24) {
		// reset the time
		lastAccessControlReset = time.Now()
		accessControl = map[string]UserAccess{}
	}

	mtx.Unlock()
}

func controlAccess(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, p httpr.Params), ctrl AccessKinds) func(ctx context.Context, w http.ResponseWriter, r *http.Request, p httpr.Params) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, p httpr.Params) {
		// get the token from the context
		token := ctx.Value("token").(string)

		// check if we need to reset the map

		mtx.Lock()
		// check if the user can still write
		if ctrl == Create {
			if ac, ok := accessControl[token]; !ok {
				newAc := UserAccess{Create: 1}
				accessControl[token] = newAc
			} else {
				if ac.Create >= maxCreate {
					writeResponse(w, http.StatusUnauthorized, "borg-api: api max create reached")
					return
				} else {
					ac.Create += 1
					accessControl[token] = ac
				}
			}
		}

		if ctrl == Update {
			if ac, ok := accessControl[token]; !ok {
				newAc := UserAccess{Update: 1}
				accessControl[token] = newAc
			} else {
				if ac.Create >= maxUpdate {
					writeResponse(w, http.StatusUnauthorized, "borg-api: api max update reached")
					return
				} else {
					ac.Create += 1
					accessControl[token] = ac
				}
			}
		}

		// just log some shit
		log.Infof("[user access control] token: %s -> %#v", token, accessControl[token])

		mtx.Unlock()

		// then call the handler
		handler(ctx, w, r, p)
	}
}

// simple helper to check if the user is auth in the application,
// if logged process the handler, or return directly
func ifAuth(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, p httpr.Params)) func(w http.ResponseWriter, r *http.Request, p httpr.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httpr.Params) {
		var token string
		if token = r.FormValue("token"); token == "" {
			if token = r.Header.Get("Authorization"); token == "" {
				if token = r.Header.Get("authorization"); token == "" {
					writeResponse(w, http.StatusUnauthorized, "borg-api: Missing access token")
					return
				}
			}
		}
		u, err := aut.GetUser(token)
		if err != nil || u == nil {
			// github may not recognize the token, return an error
			writeResponse(w, http.StatusUnauthorized, "borg-api: Invalid access token")
			return
		}
		// no errors, process the handler
		ctx := context.WithValue(context.Background(), "token", token)
		ctx = context.WithValue(ctx, "userId", u.Id)
		handler(ctx, w, r, p)
	}
}

// just redirect the user with the url to the github oauth login with the client_id
// setted in the backend
func redirectGithubAuthorize(w http.ResponseWriter, r *http.Request, p httpr.Params) {
	url := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&scope=read:org",
		*githubClientId,
	)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func githubAuth(w http.ResponseWriter, r *http.Request, p httpr.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	user, err := aut.GithubAuth(string(body))
	if err != nil {
		fmt.Fprintln(w, fmt.Sprintf("Auth failed: %v", err))
		return
	}
	bs, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(bs))
}

func getUser(ctx context.Context, w http.ResponseWriter, r *http.Request, p httpr.Params) {
	user, err := aut.GetUser(r.FormValue("token"))
	if err != nil {
		fmt.Fprintln(w, fmt.Sprintf("Getting user failed: %v", err))
		return
	}
	if user == nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	bs, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(bs))
}

func query(w http.ResponseWriter, r *http.Request, p httpr.Params) {
	size := 5
	s, err := strconv.ParseInt(r.FormValue("l"), 10, 32)
	if err == nil && s > 0 {
		size = int(s)
	}
	if size > 50 {
		size = 50
	}
	q := r.FormValue("q")
	if r.FormValue("p") == "true" {
		log.Infof("Private query with size '%v'", size)
	} else {
		log.Infof("Querying '%v' with size '%v'", q, size)
	}
	res, err := client.Search().Index("borg").Type("problem").From(0).Size(size).Query(
		elastic.NewMultiMatchQuery(q).FieldWithBoost("Title", 5.0).Field("Solutions.Body")).Do()
	if err != nil {
		panic(err)
	}
	all := []types.Problem{}
	var ttyp types.Problem
	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(types.Problem); ok {
			all = append(all, t)
		}
	}
	bs, err := json.Marshal(all)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(bs))
}

func sitemapLoop() {
	first := true
	for {
		if !first {
			time.Sleep(30 * time.Minute)
		}
		first = false
		generateSitemap()
	}
}

func generateSitemap() {
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("Sitemap generation failed: %v", r)
		}
	}()
	// this query is because we only want to show user submitted content for now - not ones scraped from somewhere else - to not piss of google
	// @TODO include ones which were changed substantially
	// @TODO this is going to get dog slow
	res, err := client.Search().Query(elastic.NewFilteredQuery(elastic.NewRegexpFilter("CreatedBy", ".{3,}"))).Size(500).Do()
	if err != nil {
		panic(err)
	}
	all := []types.Problem{}
	var ttyp types.Problem
	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(types.Problem); ok {
			all = append(all, t)
		}
	}
	items := []*sitemap.Item{}
	for _, v := range all {
		item := &sitemap.Item{
			Loc:        "http://ok-b.org/t/" + fmt.Sprintf("%v/%v", v.Id, slugify.S(v.Title)),
			LastMod:    time.Now(),
			Priority:   0.5,
			Changefreq: "daily",
		}
		items = append(items, item)
	}
	err = sitemap.SiteMap(*sm+"/sitemap.xml.gz", items)
	if err != nil {
		panic(err)
	}
	log.Info("Generated sitemap successfully")
}
