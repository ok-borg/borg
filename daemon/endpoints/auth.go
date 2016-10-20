package endpoints

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/jpillora/go-ogle-analytics"
	"github.com/ok-borg/borg/daemon/domain"
	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"gopkg.in/olivere/elastic.v3"
)

// NewEndpoints is just below the http handlers
func NewEndpoints(
	oauthCfg *oauth2.Config,
	client *elastic.Client,
	a *ga.Client,
	db *gorm.DB,
) *Endpoints {
	return &Endpoints{
		oauthCfg:  oauthCfg,
		client:    client,
		analytics: a,
		db:        db,
	}
}

// Endpoints represents all endpoints of the http server
type Endpoints struct {
	oauthCfg  *oauth2.Config
	client    *elastic.Client
	analytics *ga.Client
	db        *gorm.DB
}

func githubUserToBorgUser(user *github.User) domain.User {
	ret := domain.User{
		Id:          uuid.NewV4().String(),
		Login:       *user.Login,
		AvatarUrl:   *user.AvatarURL,
		AccountType: domain.AccountTypeGithub,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if user.Email != nil {
		ret.Email = *user.Email
	}
	if user.Name != nil {
		ret.Name = *user.Name
	}
	return ret
}

// GithubAuth exchanges a github code for a token, registers and returns a User
func (e *Endpoints) GithubAuth(code string) (*domain.User, *domain.AccessToken, error) {
	if len(code) == 0 {
		return nil, nil, errors.New("Code received is empty")
	}
	tkn, err := e.oauthCfg.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, nil, fmt.Errorf("there was an issue getting your token: %v", err)
	}
	if !tkn.Valid() {
		return nil, nil, errors.New("Reretreived invalid token")
	}
	client := github.NewClient(e.oauthCfg.Client(oauth2.NoContext, tkn))
	user, _, err := client.Users.Get("")
	if err != nil {
		return nil, nil, fmt.Errorf("error getting name: %v", err)
	}
	// here we got a github user
	// first check if a github_users row exists with this github_id.
	// if yes just save the token and associated it to the borg user linked to the github user
	// if no, create a borg_users from the github users, associated both in a github_users row
	// and finally create the access_token in db.
	ghUserDao := domain.NewGithubUserDao(e.db)

	var borgUser domain.User

	ghUser, err := ghUserDao.GetByGithubId(fmt.Sprintf("%v", *user.ID))
	if err != nil {
		// github user do not exist
		// so the borg user cannot exists too
		// first create it
		newUser := githubUserToBorgUser(user)
		userDao := domain.NewUserDao(e.db)
		if err := userDao.Create(newUser); err != nil {
			return nil, nil, fmt.Errorf("error creating new user %s", err.Error())
		}

		// create the github user to link to our new borg user
		newGithubUser := domain.GithubUser{
			Id:         uuid.NewV4().String(),
			GithubId:   strconv.Itoa(*user.ID),
			BorgUserId: newUser.Id,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if err := ghUserDao.Create(newGithubUser); err != nil {
			return nil, nil, fmt.Errorf("error creating new github user %s", err.Error())
		}

		borgUser = newUser

	} else {
		var err error
		userDao := domain.NewUserDao(e.db)
		borgUser, err = userDao.GetById(ghUser.BorgUserId)
		if err != nil {
			return nil, nil, fmt.Errorf("error getting user %s", err.Error())
		}
	}

	// then just need to try to get the access token
	tokenDao := domain.NewAccessTokenDao(e.db)
	token, err := tokenDao.GetByToken(tkn.AccessToken)
	if err != nil {
		// token do not exist in db, just create it
		token = domain.AccessToken{
			Id:        uuid.NewV4().String(),
			Token:     tkn.AccessToken,
			UserId:    borgUser.Id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := tokenDao.Create(token); err != nil {
			return nil, nil, fmt.Errorf("error creating token for user: %s", err.Error())
		}
	}

	return &borgUser, &token, nil
}

// GetUser by token
func (e *Endpoints) GetUser(token string) (*domain.User, error) {
	// first get token
	tokenDao := domain.NewAccessTokenDao(e.db)
	t, err := tokenDao.GetByToken(token)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("token (%s) is associated to no users", token))
	}
	userDao := domain.NewUserDao(e.db)
	u, _ := userDao.GetById(t.UserId)
	return &u, nil

}
