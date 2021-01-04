package authtoken

import "time"

type session interface {
	IsExpired() bool
	GetUsername() string
}

type sessionImpl struct {
	expiryTime time.Time
	username   string
	authToken  string
}

func (s *sessionImpl) IsExpired() bool {
	return time.Now().After(s.expiryTime)
}

func (s *sessionImpl) GetUsername() string {
	return s.username
}
