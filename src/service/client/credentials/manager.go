package credentials

import (
	"fmt"
	"github.com/kllla/web/src/config"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type Manager interface {
	CreateCredentials(credentials *Credentials) error
	DeleteCredentials(credentials *Credentials) error
	IsCredentialsValid(credentials *Credentials) bool
	Close() error
}

type managerImpl struct {
	credentialDao Dao
}

var DefaultManager = newDefaultManager()
var TestManager = newTestManager()

func newDefaultManager() Manager {
	return &managerImpl{credentialDao: NewDao(config.DefaultConfig)}
}

func newTestManager() Manager {
	return &managerImpl{credentialDao: NewDao(config.TestConfig)}
}

// IsCredentialsValid checks if the user is registered and details match those stored server side
func (m *managerImpl) IsCredentialsValid(credentials *Credentials) bool {
	usrCrd := m.credentialDao.GetCredentialsForUsername(credentials.Username)
	found := len(usrCrd) > 0
	if found {
		err := bcrypt.CompareHashAndPassword([]byte(usrCrd[0].PasswordHash), []byte(credentials.Password))
		if err != nil {
			return false
		}
		return true
	}
	return false
}

// CreateCredentials adds user to  store
func (m *managerImpl) CreateCredentials(credentials *Credentials) error {
	log.Printf("Creating: %s", credentials.Username)
	return m.credentialDao.CreateCredentials(credentials)
}

func (m *managerImpl) DeleteCredentials(credentials *Credentials) error {
	if !m.IsCredentialsValid(credentials) {
		return fmt.Errorf("credentials invalid failed to delete")
	}
	log.Printf("Deleting: %s", credentials.Username)
	return m.credentialDao.DeleteCredentials(credentials)
}

func (m *managerImpl) Close() error {
	return m.credentialDao.Close()
}

func (m *managerImpl) getCredentialsFromRequest(r *http.Request) *Credentials {
	return NewCredentials(r.BasicAuth())
}
