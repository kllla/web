package credentials

import (
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string
	Password string
}

type NoPassCredentials struct {
	Username     string
	PasswordHash string
}

func NewCredentials(username string, password string, ok bool) *Credentials {
	if !ok {
		return nil
	}
	return &Credentials{
		Username: username,
		Password: password,
	}
}

func (c *Credentials) ToNoPassCredentials() *NoPassCredentials {
	return &NoPassCredentials{
		Username:     c.Username,
		PasswordHash: c.getSaltedHashedPassword(),
	}

}

func (c *Credentials) getSaltedHashedPassword() string {
	// Should never err as using bcrypyt.DefaultCost
	hash, _ := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	return string(hash)
}
