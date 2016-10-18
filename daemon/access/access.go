package access

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/cihub/seelog"
	"github.com/jinzhu/gorm"
	httpr "github.com/julienschmidt/httprouter"
	"github.com/ok-borg/borg/daemon/domain"
	"golang.org/x/net/context"
)

type AccessKinds int

type UserAccess struct {
	Update int
	Create int
}

// FIXME(jeremy): should be in config
// maximum access for write and updates in 24 hours
const (
	maxCreate = 100
	maxUpdate = 50
)

// acces kings
const (
	Create AccessKinds = iota
	Update
)

var (
	accessControl          map[string]UserAccess
	mtx                    = &sync.Mutex{}
	lastAccessControlReset = time.Now()
)

func init() {
	accessControl = map[string]UserAccess{}
}

func updateTimer() {
	mtx.Lock()
	if time.Since(lastAccessControlReset) >= (time.Hour * 24) {
		lastAccessControlReset = time.Now()
		accessControl = map[string]UserAccess{}
	}
	mtx.Unlock()
}

func Control(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, p httpr.Params), ctrl AccessKinds) func(ctx context.Context, w http.ResponseWriter, r *http.Request, p httpr.Params) {
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
func IfAuth(db *gorm.DB, handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, p httpr.Params)) func(w http.ResponseWriter, r *http.Request, p httpr.Params) {
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

		accessTokenDao := domain.NewAccessTokenDao(db)
		at, err := accessTokenDao.GetByToken(token)
		if err != nil {
			writeResponse(w, http.StatusUnauthorized, "borg-api: Invalid access token")
			return
		}
		// get or create it in mysql
		userDao := domain.NewUserDao(db)
		user, err := userDao.GetById(at.UserId)
		if err != nil {
			writeResponse(w, http.StatusUnauthorized, "borg-api: Invalid access token")
			return
		}

		// no errors, process the handler
		ctx := context.WithValue(context.Background(), "token", token)
		ctx = context.WithValue(ctx, "userId", user.Id)
		ctx = context.WithValue(ctx, "user", user)
		handler(ctx, w, r, p)
	}
}

// FIXME: this is duplicated in main.go
func writeResponse(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(body)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `%v`, body)
}
