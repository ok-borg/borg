package endpoints

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/google/go-github/github"
	"github.com/jpillora/go-ogle-analytics"
	"golang.org/x/oauth2"
	"gopkg.in/olivere/elastic.v3"
)

// NewEndpoints is just below the http handlers
func NewEndpoints(oauthCfg *oauth2.Config, client *elastic.Client, a *ga.Client) *Endpoints {
	return &Endpoints{
		oauthCfg:  oauthCfg,
		client:    client,
		analytics: a,
	}
}

// Endpoints represents all endpoints of the http server
type Endpoints struct {
	oauthCfg  *oauth2.Config
	client    *elastic.Client
	analytics *ga.Client
}

// GithubAuth exchanges a github code for a token, registers and returns a User
func (e *Endpoints) GithubAuth(code string) (*User, error) {
	if len(code) == 0 {
		return nil, errors.New("Code received is empty")
	}
	tkn, err := e.oauthCfg.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("there was an issue getting your token: %v", err)
	}
	if !tkn.Valid() {
		return nil, errors.New("Reretreived invalid token")
	}
	client := github.NewClient(e.oauthCfg.Client(oauth2.NoContext, tkn))
	user, _, err := client.Users.Get("")
	if err != nil {
		return nil, fmt.Errorf("error getting name: %v", err)
	}
	usr, err := toUser(user)
	if err != nil {
		return nil, fmt.Errorf("error converting user: %v", err)
	}
	usr.Token = tkn.AccessToken
	// we just set the user every time for now. reuse github id. save token next to it. identify
	// user by querying users with that token.
	err = e.setUser(*usr)
	if err != nil {
		return nil, fmt.Errorf("error converting user: %v", err)
	}
	return usr, nil
}

// GetUser by token
func (e *Endpoints) GetUser(token string) (*User, error) {
	return e.readUser("Token", token)
}

// User represents a borg user
type User struct {
	Id       string
	Login    string
	Email    string
	Name     string
	SourceId string
	Token    string
}

func toUser(user *github.User) (*User, error) {
	id := fmt.Sprintf("%v", *user.ID)
	ret := &User{
		Id:       id,
		Login:    *user.Login,
		SourceId: id,
	}
	if user.Email != nil {
		ret.Email = *user.Email
	}
	if user.Name != nil {
		ret.Name = *user.Name
	}
	return ret, nil
}

func (e *Endpoints) readUser(field, equalsTo string) (*User, error) {
	termQuery := elastic.NewTermQuery(field, equalsTo)
	res, err := e.client.Search().Index("borg").Type("user").Query(termQuery).From(0).Size(2).Do()
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
		return nil, nil
	case len(users) > 1:
		return nil, fmt.Errorf("Multiple users found with %v %v ", field, equalsTo)
	}
	return &users[0], nil
}

// register or update the token
func (e *Endpoints) setUser(user User) error {
	_, err := e.client.Index().
		Index("borg").
		Type("user").
		Id(user.Id).
		BodyJson(user).
		Refresh(true).
		Do()
	return err
}
