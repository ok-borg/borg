package auth

import (
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	//"github.com/ventu-io/go-shortid"
	"golang.org/x/oauth2"
	"gopkg.in/olivere/elastic.v2"

	"reflect"
)

func NewAuth(oauthCfg *oauth2.Config, client *elastic.Client) *Auth {
	return &Auth{
		oauthCfg: oauthCfg,
		client:   client,
	}
}

type Auth struct {
	oauthCfg *oauth2.Config
	client   *elastic.Client
}

func (auth *Auth) GithubAuth(code string) (*User, error) {
	if len(code) == 0 {
		return nil, errors.New("Code received is empty")
	}
	tkn, err := auth.oauthCfg.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("there was an issue getting your token: %v", err))
	}
	if !tkn.Valid() {
		return nil, errors.New("Reretreived invalid token")
	}
	client := github.NewClient(auth.oauthCfg.Client(oauth2.NoContext, tkn))
	user, _, err := client.Users.Get("")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting name: %v", err))
	}
	usr, err := toUser(user)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error converting user: %v", err))
	}
	usr.Token = tkn.AccessToken
	// we just set the user every time for now. reuse github id. save token next to it. identify
	// user by querying users with that token.
	err = auth.setUser(*usr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error converting user: %v", err))
	}
	return usr, nil
}

func (auth *Auth) GetUser(token string) (*User, error) {
	return auth.readUser("Token", token)
}

type User struct {
	Id       string
	Email    string
	Name     string
	SourceId string
	Token    string
}

func toUser(user *github.User) (*User, error) {
	var email string
	switch {
	case user.Email != nil:
		email = *user.Email
	case user.Name == nil:
		return nil, errors.New("User has no name") // hopefully this is impossible
	case user.Name != nil:
		email = *user.Name // FIXME(crufter): this is just a quick temporary fix - not everyone has a public email at github
	}
	id := fmt.Sprintf("%v", *user.ID)
	ret := &User{
		Id:       id,
		Email:    email,
		Name:     *user.Name,
		SourceId: id,
	}
	return ret, nil
}

func (auth *Auth) readUser(field, equalsTo string) (*User, error) {
	termQuery := elastic.NewTermQuery(field, equalsTo)
	res, err := auth.client.Search().Index("borg").Type("user").Query(termQuery).From(0).Size(2).Do()
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
		return nil, errors.New(fmt.Sprintf("Multiple users found with %v %v ", field, equalsTo))
	}
	return &users[0], nil
}

// register or update the token
func (auth *Auth) setUser(user User) error {
	_, err := auth.client.Index().
		Index("borg").
		Type("user").
		Id(user.Id).
		BodyJson(user).
		Refresh(true).
		Do()
	return err
}
