package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/cihub/seelog"
	"github.com/crufter/borg/daemon/access"
	"github.com/crufter/borg/daemon/endpoints"
	"github.com/crufter/borg/daemon/sitemap"
	"github.com/crufter/borg/types"
	"github.com/jpillora/go-ogle-analytics"
	httpr "github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"golang.org/x/oauth2"
	"gopkg.in/olivere/elastic.v3"
)

const (
	githubTokenURL = "https://github.com/login/oauth/access_token"
)

var (
	esAddr             = flag.String("esaddr", "127.0.0.1:9200", "Elastic Search address")
	githubClientId     = flag.String("github-client-id", "", "Github oauth client id")
	githubClientSecret = flag.String("github-client-secret", "", "Github client secret")
	sm                 = flag.String("sitemap", "", "Sitemap location. Leave empty if you don't want a sitemap to be generated")
	analytics          = flag.String("analytics", "", "Analytics tracking id")
)

var (
	client          *elastic.Client
	analyticsClient *ga.Client
	ep              *endpoints.Endpoints
)

type Logger struct{}

func (l Logger) Printf(str string, i ...interface{}) {
	fmt.Println(fmt.Sprintf(str, i...))
}

func init() {
	flag.Parse()
	cl, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(fmt.Sprintf("http://%v", *esAddr)))
	if err != nil {
		panic(err)
	}
	client = cl
	if len(*analytics) > 0 {
		acl, err := ga.NewClient(*analytics)
		if err != nil {
			log.Errorf("Failed to acquire analytics client id: %v", err)
		}
		analyticsClient = acl
	}
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
	ep = endpoints.NewEndpoints(oauthCfg, client, analyticsClient)
	r := httpr.New()
	if len(*sm) > 0 {
		go sitemapLoop(*sm, client)
	}

	r.GET("/v1/redirect/github/authorize", redirectGithubAuthorize)
	r.GET("/v1/query", q)
	r.POST("/v1/auth/github", githubAuth)

	// authenticated endpoints
	r.GET("/v1/user", access.IfAuth(client, getUser))

	// snippets
	r.GET("/v1/p/:id", getSnippet)
	r.GET("/v1/latest", getLatestSnippets)
	r.POST("/v1/p", access.IfAuth(client, access.Control(createSnippet, access.Create)))
	//r.DELETE("/v1/p/:id", access.IfAuth(deleteSnippet))
	r.PUT("/v1/p", access.IfAuth(client, access.Control(updateSnippet, access.Update)))
	r.POST("/v1/worked", access.IfAuth(client, snippetWorked))

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
	user, err := ep.GithubAuth(string(body))
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
	user, err := ep.GetUser(r.FormValue("token"))
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

func q(w http.ResponseWriter, r *http.Request, p httpr.Params) {
	size := 5
	s, err := strconv.ParseInt(r.FormValue("l"), 10, 32)
	if err == nil && s > 0 {
		size = int(s)
	}
	res, err := ep.Query(r.FormValue("q"), size, r.FormValue("p") == "true")
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err.Error())
	}
	bs, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(bs))
}

func sitemapLoop(path string, client *elastic.Client) {
	first := true
	for {
		if !first {
			time.Sleep(30 * time.Minute)
		}
		first = false
		sitemap.GenerateSitemap(path, client)
	}
}

func getLatestSnippets(w http.ResponseWriter, r *http.Request, p httpr.Params) {
	res, err := ep.GetLatestSnippets()
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err.Error())
	}
	bs, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	writeResponse(w, http.StatusOK, string(bs))
}

func createSnippet(ctx context.Context, w http.ResponseWriter, r *http.Request, p httpr.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "borg-api: unable to read body")
		return
	}
	var snipp types.Problem
	if err := json.Unmarshal(body, &snipp); err != nil {
		log.Errorf("Invalid snippet, %s, input was %s", err.Error(), string(body))
		writeResponse(w, http.StatusBadRequest, "borg-api: Invalid snippet")
		return
	}
	err = ep.CreateSnippet(snipp, ctx.Value("userId").(string))
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "borg-api: unable to unmarshal snippet")
		return
	}
	writeJsonResponse(w, http.StatusOK, snipp)
}

func updateSnippet(ctx context.Context, w http.ResponseWriter, r *http.Request, p httpr.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "borg-api: unable to read body")
		return
	}
	var snipp types.Problem
	if err := json.Unmarshal(body, &snipp); err != nil {
		log.Errorf("[updateSnippet] invalid snippet, %s, input was %s", err.Error(), string(body))
		writeResponse(w, http.StatusBadRequest, "borg-api: Invalid snippet")
		return
	}
	err = ep.UpdateSnippet(snipp, ctx.Value("userId").(string))
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "borg-api: error")
		return
	}
	writeResponse(w, http.StatusOK, "{}")
}

func snippetWorked(ctx context.Context, w http.ResponseWriter, r *http.Request, p httpr.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "borg-api: unable to read body")
		return
	}
	s := struct {
		Query string
		Id    string
	}{}
	if err := json.Unmarshal(body, &s); err != nil {
		log.Errorf("[updateSnippet] invalid worked request, %s, input was %s", err.Error(), string(body))
		writeResponse(w, http.StatusBadRequest, "borg-api: Invalid worked request")
		return
	}
	err = ep.Worked(s.Id, s.Query)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "borg-api: error: "+err.Error())
		return
	}
	writeResponse(w, http.StatusOK, "{}")
}

func getSnippet(w http.ResponseWriter, r *http.Request, p httpr.Params) {
	id := p.ByName("id")
	if len(id) == 0 {
		writeResponse(w, http.StatusBadRequest, "borg-api: Missing id url parameter")
		return
	}
	snipp, err := ep.GetSnippet(id)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "borg-api: Failed to get snippet")
		return
	}
	if snipp == nil {
		writeResponse(w, http.StatusNotFound, "borg-api: snippet not found")
		return
	}
	bs, _ := json.Marshal(snipp)
	writeResponse(w, http.StatusOK, string(bs))
}
