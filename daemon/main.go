package main

import (
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/crufter/borg/daemon/auth"
	"github.com/crufter/borg/types"
	httpr "github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"golang.org/x/oauth2"
	"gopkg.in/olivere/elastic.v2"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
)

const (
	githubTokenURL = "https://github.com/login/oauth/access_token"
)

var (
	esAddr             = flag.String("esaddr", "127.0.0.1:9200", "Elastic Search address")
	githubClientId     = flag.String("github-client-id", "", "Github oauth client id")
	githubClientSecret = flag.String("github-client-secret", "", "Github client secret")
)

var (
	client *elastic.Client
	aut    *auth.Auth
)

type Logger struct{}

func (l *Logger) Printf(str string, i ...interface{}) {
	fmt.Println(fmt.Sprintf(str, i...))
}

func init() {
	flag.Parse()
	cl, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(fmt.Sprintf("http://%v", *esAddr)))
	if err != nil {
		panic(err)
	}
	client = cl
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
	r.GET("/v1/query", query)
	r.POST("/v1/auth/github", githubAuth)
	handler := cors.Default().Handler(r)
	log.Info("Starting http server")
	log.Critical(http.ListenAndServe(fmt.Sprintf(":%v", 9992), handler))
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
