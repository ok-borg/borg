package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/olivere/elastic.v3"

	"golang.org/x/net/context"

	log "github.com/cihub/seelog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jpillora/go-ogle-analytics"
	httpr "github.com/julienschmidt/httprouter"
	"github.com/ok-borg/borg/daemon/access"
	"github.com/ok-borg/borg/daemon/conf"
	"github.com/ok-borg/borg/daemon/domain"
	"github.com/ok-borg/borg/daemon/endpoints"
	"github.com/ok-borg/borg/daemon/sitemap"
	"github.com/ok-borg/borg/types"
	"github.com/rs/cors"
	"golang.org/x/oauth2"
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
	sqlAddr            = flag.String("sqladdr", "127.0.0.1:3306", "Mysql address")
	sqlIds             = flag.String("sqlids", "root:root", "Mysql identifier")
)

var (
	client          *elastic.Client
	analyticsClient *ga.Client
	ep              *endpoints.Endpoints
	db              *gorm.DB
)

func initWithConfFile() {
	// let's assume the conf file is in the same place than the benary.
	c, err := ioutil.ReadFile(".borg.conf.json")
	if err != nil {
		// just return, there is probably not configuration file
		return
	}

	var conf conf.Conf
	if err := json.Unmarshal(c, &conf); err != nil {
		panic(fmt.Sprintf("[initWithConfFile] invalid config format: %s", err.Error()))
	}

	if conf.EsAddr != "" {
		*esAddr = conf.EsAddr
	}
	if conf.Github.ClientId != "" {
		*githubClientId = conf.Github.ClientId
	}
	if conf.Github.ClientSecret != "" {
		*githubClientSecret = conf.Github.ClientSecret
	}
	if conf.Sitemap != "" {
		*sm = conf.Sitemap
	}
	if conf.Analytics != "" {
		*analytics = conf.Analytics
	}
	if conf.Mysql.Addr != "" {
		*sqlAddr = conf.Mysql.Addr
	}
	if conf.Mysql.Ids != "" {
		*sqlIds = conf.Mysql.Ids
	}
}

func initLogger() {
	logger, _ := log.LoggerFromConfigAsString(conf.Seelog)
	log.ReplaceLogger(logger)
}

func init() {
	// read config file before if it exists, so we can replaces the var that was set with the cmdline
	// the cmdline is allowed to overwrite the config file.
	initWithConfFile()
	flag.Parse()
	initLogger()

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

	// init mysql
	var err error
	dsn := fmt.Sprintf("%s@tcp(%s)/borg?parseTime=True", *sqlIds, *sqlAddr)
	if db, err = gorm.Open("mysql", dsn); err != nil {
		panic(fmt.Sprintf("[init] unable to initialize gorm: %s", err.Error()))
	}
	defer db.Close()

	// decl routes

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
	r.POST("/v1/slack", slackCommand)

	// organizations
	r.POST("/v1/organizations", access.IfAuth(client, createOrganization))
	r.GET("/v1/organizations", access.IfAuth(client, listUserOrganizations))

	// not rest at all but who cares ?
	r.POST("/v1/organizations/leave/:id", access.IfAuth(client, leaveOrganization))
	r.POST("/v1/organizations/expel/:oid/user/id/:uid",
		access.IfAuth(client, expelUserFromOrganization))
	r.POST("/v1/organizations/admins/:oid/user/id/:uid",
		access.IfAuth(client, grantAdminRightToUser))

	// organizations-join-links
	// this is only allowed for the organization admin
	r.POST("/v1/organization-join-links", access.IfAuth(client, createOrganizationJoinLink))
	r.DELETE("/v1/organization-join-links/id/:id", access.IfAuth(client, deleteOrganizationJoinLink))
	// get a join link for a specific organization
	// this is allowed only by the organization admin in order to share it again, or delete it.
	r.GET("/v1/organization-join-links/organizations/:id",
		access.IfAuth(client, getOrganizationJoinLinkByOrganizationId))
	// get a join link from a join-link id.
	r.GET("/v1/organization-join-links/id/:id", access.IfAuth(client, getOrganizationJoinLink))
	// accept join link
	// not restful at all, but pretty to read
	r.POST("/v1/join/:id", access.IfAuth(client, joinOrganization))

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
	err = ep.CreateSnippet(&snipp, ctx.Value("userId").(string))
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
	err = ep.UpdateSnippet(&snipp, ctx.Value("userId").(string))
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

func readJsonBody(r *http.Request, expectedBody interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.New("unable to read body")
	}
	if err := json.Unmarshal(body, expectedBody); err != nil {
		log.Errorf(
			"[readJsonBody] invalid request, %s, input was %s",
			err.Error(), string(body))
		return errors.New("invalid json body format")
	}
	return nil
}

func getUserByAccessToken(ctx context.Context) (domain.User, error) {
	// get user in elastic
	u, _ := ep.GetUser(ctx.Value("token").(string))
	// get or create it in mysql
	userDao := domain.NewUserDao(db)
	return userDao.GetOrCreateFromRaw(u.Login, u.Email, u.Id)
}

func createOrganization(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	p httpr.Params) {

	// first unmarshal body
	expectedBody := struct{ Name string }{}
	if err := readJsonBody(r, &expectedBody); err != nil {
		writeResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("borg-api: %s", err.Error()))
		return
	}

	if u, err := getUserByAccessToken(ctx); err != nil {
		// handle shit here
	} else {
		// lets create an org
		if o, err := ep.CreateOrganization(db, u.Id, expectedBody.Name); err != nil {
			writeResponse(w, http.StatusInternalServerError, "borg-api: create organization error: "+err.Error())
			return
		} else {
			writeJsonResponse(w, http.StatusOK, o)
		}
	}
}

// create a new organization join link.
// only an administrator of an organization can execute this action
func createOrganizationJoinLink(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	p httpr.Params) {

	// first unmarshal body
	expectedBody := struct {
		OrganizationId string
		Ttl            int64
	}{}
	if err := readJsonBody(r, &expectedBody); err != nil {
		writeResponse(w, http.StatusInternalServerError,
			fmt.Sprintf("borg-api: %s", err.Error()))
		return
	}

	// check mandatory fields
	if expectedBody.OrganizationId == "" || expectedBody.Ttl <= 0 {
		log.Errorf(
			"[createOrganizationJoinLink] invalid createOrganizationjoinlink body")
		writeResponse(w, http.StatusBadRequest, "borg-api: invalid body")
		return
	}

	if u, err := getUserByAccessToken(ctx); err != nil {
		// handle shit here
	} else {
		// ceate the organizartion Join Link
		if o, err := ep.CreateOrganizationJoinLink(db, u.Id, expectedBody.OrganizationId, expectedBody.Ttl); err != nil {
			writeResponse(w, http.StatusInternalServerError,
				"borg-api: create organization join link error: "+err.Error())
			return
		} else {
			writeJsonResponse(w, http.StatusOK, o)
		}

	}
}

// delete an existing link
// same as previously
func deleteOrganizationJoinLink(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	p httpr.Params) {
	id := p.ByName("id")
	if len(id) == 0 {
		writeResponse(w, http.StatusBadRequest, "borg-api: Missing id url parameter")
		return
	}

	if u, err := getUserByAccessToken(ctx); err != nil {
		// handle shit here
	} else {
		// delete the organizartion Join Link
		if err := ep.DeleteOrganizationJoinLink(db, u.Id, id); err != nil {
			writeResponse(w, http.StatusInternalServerError,
				"borg-api: delete organization join link error: "+err.Error())
			return
		}
		writeResponse(w, http.StatusOK, "")
	}
}

// by id
// get an existing link in order to consult the time left for the
// join link, or delete it, or get the the organizastion link for invited
// users to display orgs infos
func getOrganizationJoinLink(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	p httpr.Params) {
	id := p.ByName("id")
	if len(id) == 0 {
		writeResponse(w, http.StatusBadRequest, "borg-api: Missing id url parameter")
		return
	}

	// get the organization join link
	// no need of user id or anythin
	if ojl, err := ep.GetOrganizationJoinLink(db, id); err != nil {
		writeResponse(w, http.StatusInternalServerError,
			"borg-api: get organization join link error: "+err.Error())
		return
	} else {
		writeJsonResponse(w, http.StatusOK, ojl)
	}
}

// get join link for a given organization
// will work only for an admin in order to manage this join-link
func getOrganizationJoinLinkByOrganizationId(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	p httpr.Params) {
	id := p.ByName("id")
	if len(id) == 0 {
		writeResponse(w, http.StatusBadRequest, "borg-api: Missing id url parameter")
		return
	}

	if u, err := getUserByAccessToken(ctx); err != nil {
		// handle shit here
	} else {
		// ceate the organizartion Join Link
		if ojl, err := ep.GetOrganizationJoinLinkForOrganization(db, u.Id, id); err != nil {
			writeResponse(w, http.StatusInternalServerError,
				"borg-api: get organization join link error: "+err.Error())
			return
		} else {
			writeJsonResponse(w, http.StatusOK, ojl)
		}
	}
}

// join an organization.
// if join link is not expired.
func joinOrganization(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	p httpr.Params) {
	id := p.ByName("id")
	if len(id) == 0 {
		writeResponse(w, http.StatusBadRequest, "borg-api: Missing id url parameter")
		return
	}

	if u, err := getUserByAccessToken(ctx); err != nil {
		// handle shit here
	} else {
		// ceate the organizartion Join Link
		if err := ep.JoinOrganization(db, u.Id, id); err != nil {
			writeResponse(w, http.StatusInternalServerError,
				"borg-api: cannot join organization: "+err.Error())
			return
		} else {
			writeJsonResponse(w, http.StatusOK, "")
		}
	}
}

// list user organization
func listUserOrganizations(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	p httpr.Params) {
	if u, err := getUserByAccessToken(ctx); err != nil {
		// handle shit here
	} else {
		// ceate the organizartion Join Link
		if orgz, err := ep.ListUserOrganizations(db, u.Id); err != nil {
			writeResponse(w, http.StatusInternalServerError,
				"borg-api: list user organizations error: "+err.Error())
			return
		} else {
			writeJsonResponse(w, http.StatusOK, orgz)
		}
	}

}

// leave an organization,
// you cannot leave an organization if you are the only admin for it
func leaveOrganization(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	p httpr.Params,
) {
	organizationId := p.ByName("id")
	if len(organizationId) == 0 {
		writeResponse(w, http.StatusBadRequest, "borg-api: Missing id url parameter")
		return
	}

	if u, err := getUserByAccessToken(ctx); err != nil {
		// handle shit here
	} else {
		// ceate the organizartion Join Link
		if err := ep.LeaveOrganization(db, u.Id, organizationId); err != nil {
			writeResponse(w, http.StatusInternalServerError,
				"borg-api: cannot leave organization: "+err.Error())
			return
		} else {
			writeJsonResponse(w, http.StatusNoContent, "")
		}
	}

}

// expel an user from an organization,
// you can only do this if you are admin of the organization from
// where you want to expel someone
func expelUserFromOrganization(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	p httpr.Params,
) {
	organizationId := p.ByName("oid")
	if len(organizationId) == 0 {
		writeResponse(w, http.StatusBadRequest, "borg-api: Missing organizationId url parameter")
		return
	}
	userId := p.ByName("uid")
	if len(userId) == 0 {
		writeResponse(w, http.StatusBadRequest, "borg-api: Missing userId url parameter")
		return
	}

	if u, err := getUserByAccessToken(ctx); err != nil {
		// handle shit here
	} else {
		// ceate the organizartion Join Link
		if err := ep.ExpelUserFromOrganization(db, u.Id, userId, organizationId); err != nil {
			writeResponse(w, http.StatusInternalServerError,
				"borg-api: cannot expel from organization: "+err.Error())
			return
		} else {
			writeJsonResponse(w, http.StatusNoContent, "")
		}
	}
}

func grantAdminRightToUser(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	p httpr.Params,
) {
	organizationId := p.ByName("oid")
	if len(organizationId) == 0 {
		writeResponse(w, http.StatusBadRequest, "borg-api: Missing organizationId url parameter")
		return
	}
	userId := p.ByName("uid")
	if len(userId) == 0 {
		writeResponse(w, http.StatusBadRequest, "borg-api: Missing userId url parameter")
		return
	}

	if u, err := getUserByAccessToken(ctx); err != nil {
		// handle shit here
	} else {
		// ceate the organizartion Join Link
		if err := ep.GrantAdminRightToUser(db, u.Id, userId, organizationId); err != nil {
			writeResponse(w, http.StatusInternalServerError,
				"borg-api: cannot expel from organization: "+err.Error())
			return
		} else {
			writeJsonResponse(w, http.StatusNoContent, "")
		}
	}

}

func slackCommand(w http.ResponseWriter, r *http.Request, p httpr.Params) {
	if err := r.ParseForm(); err != nil {
		writeResponse(w, http.StatusInternalServerError, "Something wrong happened, please try again later.")
		return
	}
	if res, err := ep.Slack(r.FormValue("text")); err != nil {
		writeResponse(w, http.StatusInternalServerError, "Something wrong happened, please try again later.")
		return
	} else {
		writeResponse(w, http.StatusOK, res)
	}
}
