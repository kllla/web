package authtoken

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kllla/web/src/service/client/credentials"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
)

const (
	cookieKey string = "auth"
)

var sessionKey = "SESSION_KEY"

var DefaultManager = newManager(sessionKey)

type Service interface {
	AuthenticateSession(w http.ResponseWriter, r *http.Request, creds *credentials.Credentials) bool
	UnAuthenticateSession(w http.ResponseWriter, r *http.Request)
	IsSessionAuthenticated(w http.ResponseWriter, r *http.Request) bool
	GetUsernameForSession(w http.ResponseWriter, r *http.Request) string
}

// impl is for storing http user sessions
type impl struct {
	cookieKey      string
	cookieStore    *sessions.CookieStore
	activeSessions map[string]session
	defaultExpiry  time.Duration
}

func (sh *impl) GetUsernameForSession(w http.ResponseWriter, r *http.Request) string {
	fmt.Print(fmt.Sprintf("sh.activeSessions %v", sh.activeSessions), fmt.Sprintf("Getting Session Cookie for mck: %s vs ck: %s ", sh.cookieKey, cookieKey))
	cookies, err := sh.cookieStore.Get(r, sh.cookieKey)
	if err == nil {
		if sesKey, cok := cookies.Values[cookieKey]; cok {
			if ses, ok := sh.activeSessions[fmt.Sprintf("%v", sesKey)]; ok {

				if !ses.(session).IsExpired() {
					return ses.GetUsername()
				}
			}
		}
	}
	return ""
}

func (sh *impl) IsSessionAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	cookies, err := sh.cookieStore.Get(r, sh.cookieKey)
	if err == nil {
		if sesKey, ok := cookies.Values[cookieKey]; ok {
			if ses, ok := sh.activeSessions[fmt.Sprintf("%v", sesKey)]; ok {
				if !ses.(session).IsExpired() {
					return true
				}
			}
		}
	}
	return false
}

// NewManager returns a new manager with key provided
func newManager(sessionkey string) Service {
	return &impl{
		cookieStore:    sessions.NewCookieStore([]byte(sessionkey)),
		activeSessions: make(map[string]session, 0),
		cookieKey:      sessionKey,
		defaultExpiry:  time.Minute * 30,
	}
}

// AuthenticateUser adds user to session cookieStore
func (sh *impl) AuthenticateSession(w http.ResponseWriter, r *http.Request, creds *credentials.Credentials) bool {
	cookies, _ := sh.cookieStore.Get(r, sh.cookieKey)
	sessionAuthToken := sh.createAuthToken()
	authedSession := &sessionImpl{
		username:   creds.Username,
		expiryTime: time.Now().Add(sh.defaultExpiry),
		authToken:  sessionAuthToken,
	}
	sh.activeSessions[sessionAuthToken] = authedSession
	cookies.Values[cookieKey] = sessionAuthToken
	cookies.Save(r, w)
	return true
}

// DeauthenticateUser removes user from session cookieStore
func (sh *impl) UnAuthenticateSession(w http.ResponseWriter, r *http.Request) {
	cookies, _ := sh.cookieStore.Get(r, sh.cookieKey)
	sessionToken := cookies.Values[cookieKey]
	if sessionToken != nil {
		delete(sh.activeSessions, fmt.Sprintf("%v", sessionToken))
		cookies.Values[cookieKey] = fmt.Sprintf("expired %s", time.Now())
		cookies.Save(r, w)
	}
}

func (sh *impl) createAuthToken() string {
	return uuid.New().String()
}
