package client

import (
	"fmt"
	"github.com/kllla/web/src/service/client/authtoken"
	"github.com/kllla/web/src/service/client/credentials"
	"net/http"
)

type Service interface {
	GetCredentials(w http.ResponseWriter, r *http.Request) *credentials.Credentials
	CreateCredentials(w http.ResponseWriter, r *http.Request) bool
	DeleteCredentials(w http.ResponseWriter, r *http.Request) bool
	VerifyCredentialsAndAuthenticate(w http.ResponseWriter, r *http.Request) bool
	AuthenticationCheck(w http.ResponseWriter, r *http.Request) bool
	UnAuthentication(w http.ResponseWriter, r *http.Request) bool
	GetSessionUsername(w http.ResponseWriter, r *http.Request) string
}

type impl struct {
	credentialsManager credentials.Manager
	authTokenManager   authtoken.Service
}

func NewService() Service {
	return &impl{
		credentialsManager: credentials.DefaultManager,
		authTokenManager:   authtoken.DefaultManager,
	}
}

func (h *impl) GetSessionUsername(w http.ResponseWriter, r *http.Request) string {
	if username := h.authTokenManager.GetUsernameForSession(w, r); username != "" {

		return username
	}
	return ""
}

func (h *impl) GetCredentials(w http.ResponseWriter, r *http.Request) *credentials.Credentials {
	creds := h.getCredentialsFromFormData(w, r)
	return creds
}

func (h impl) CreateCredentials(w http.ResponseWriter, r *http.Request) bool {
	creds := h.getCredentialsFromFormData(w, r)
	if err := h.credentialsManager.CreateCredentials(creds); err != nil {
		http.Error(w, fmt.Sprintf("failed to create credentials: %s ", err), http.StatusInternalServerError)
		return false
	}
	return true
}

func (h impl) DeleteCredentials(w http.ResponseWriter, r *http.Request) bool {
	creds := h.getCredentialsFromFormData(w, r)
	if err := h.credentialsManager.DeleteCredentials(creds); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete credentials: %s ", err), http.StatusInternalServerError)
		return false
	}
	return true
}
func (h *impl) UnAuthentication(w http.ResponseWriter, r *http.Request) bool {
	h.authTokenManager.UnAuthenticateSession(w, r)
	return true
}

func (h *impl) AuthenticationCheck(w http.ResponseWriter, r *http.Request) bool {
	if h.authTokenManager.IsSessionAuthenticated(w, r) {
		return true
	}
	return false
}

func (h impl) VerifyCredentialsAndAuthenticate(w http.ResponseWriter, r *http.Request) bool {
	creds := h.getCredentialsFromFormData(w, r)

	if h.credentialsManager.IsCredentialsValid(creds) {

		if h.authTokenManager.AuthenticateSession(w, r, creds) {

			return true
		}
	}
	return false
}

func (h impl) getCredentialsFromFormData(w http.ResponseWriter, r *http.Request) *credentials.Credentials {
	ok := false
	if err := r.ParseForm(); err == nil {
		ok = true
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	creds := credentials.NewCredentials(username, password, ok)

	if creds == nil {
		http.Error(w, "failed basic auth form", http.StatusBadRequest)
	}
	return creds
}
