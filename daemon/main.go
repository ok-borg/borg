package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/crufter/borg/types"
	"github.com/google/go-github/github"
	httpr "github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/ventu-io/go-shortid"
	"golang.org/x/oauth2"
	"gopkg.in/olivere/elastic.v2"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
)

const (
	githubAccessURL = "https://github.com/login/oauth/access_token"
)

var (
	githubClientId     = flag.String("github-client-id", "", "Github oauth client id")
	githubClientSecret = flag.String("github-client-secret", "", "Github client secret")
	esAddr             = flag.String("esaddr", "127.0.0.1:9200", "Elastic Search address")
)

var (
	client   *elastic.Client
	oauthCfg *oauth2.Config
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
	oauthCfg = &oauth2.Config{
		ClientID:     *githubClientId,
		ClientSecret: *githubClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL: githubAccessURL,
		},
		//RedirectURL: redirectUrl,
		Scopes: []string{"read:org"},
	}
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
	tkn, err := oauthCfg.Exchange(oauth2.NoContext, string(body))
	if err != nil {
		fmt.Fprintln(w, fmt.Sprintf("there was an issue getting your token: %v", err))
		return
	}
	if !tkn.Valid() {
		fmt.Fprintln(w, "retreived invalid token")
		return
	}
	client := github.NewClient(oauthCfg.Client(oauth2.NoContext, tkn))
	user, _, err := client.Users.Get("")
	if err != nil {
		fmt.Println(w, fmt.Sprintf("error getting name: %v", err))
		return
	}
	usr, err := readUser(*user.Email)
	if err != nil {
		fmt.Println(w, fmt.Sprintf("error getting user: %v", err))
		return
	}
	if usr == nil {
		usr, err = toUser(user)
		if err != nil {
			fmt.Println(w, fmt.Sprintf("error converting user: %v", err))
			return
		}
		err = registerUser(*usr)
		if err != nil {
		}
	}
	bs, err := json.Marshal(map[string]interface{}{
		"User":  user,
		"Token": tkn.AccessToken,
	})
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(bs))
}

type User struct {
	Id    string
	Email string
	Name  string
}

func toUser(user *github.User) (*User, error) {
	switch {
	case user.Email == nil:
		return nil, errors.New("User has no email")
	case user.Name == nil:
		return nil, errors.New("User has no email")
	}
	id, err := shortid.Generate()
	if err != nil {
		return nil, err
	}
	ret := &User{
		Id:    id,
		Email: *user.Email,
		Name:  *user.Name,
	}
	return ret, nil
}

func readUser(email string) (*User, error) {
	termQuery := elastic.NewTermQuery("email", email)
	res, err := client.Search().Index("borg").Type("user").Query(termQuery).From(0).Size(2).Do()
	if err != nil {
		return nil, err
	}
	var ttyp User
	users := []User{}
	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(User); ok {
			users = append(users, t)
		}
	}
	switch {
	case len(users) == 0:
		return nil, errors.New("User not found with email " + email)
	case len(users) > 1:
		return nil, errors.New("Multiple users found with email " + email)
	}
	return &users[0], nil
}

func registerUser(user User) error {
	_, err := client.Index().
		Index("borg").
		Type("user").
		Id(user.Id).
		BodyJson(user).
		Refresh(true).
		Do()
	return err
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
