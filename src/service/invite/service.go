package invite

import (
	"fmt"
	"github.com/kllla/web/src/config"
	"github.com/kllla/web/src/service/id"
	"net/http"
	"time"
)

type Service interface {
	CreateInvite(w http.ResponseWriter, r *http.Request, createdBy string) (*Invite, error)
	DeleteInvite(id string) error
	GetInviteIfIDisValid(id string) (bool, *Invite)
	GetAllInvitesCreatedBy(by string) []*Invite
}

type impl struct {
	urlDao        *Dao
	defaultExpiry time.Duration
}

func (s *impl) DeleteInvite(id string) error {
	return s.urlDao.DeleteInviteByID(id)
}

func (s *impl) GetAllInvitesCreatedBy(createdBy string) []*Invite {
	return s.urlDao.GetInvitesCreatedBy(createdBy)
}

const aWeek = (time.Hour * 24) * 7

func NewService() Service {
	return &impl{
		urlDao:        NewDao(config.DefaultConfig),
		defaultExpiry: aWeek,
	}
}

func (s *impl) CreateInvite(w http.ResponseWriter, r *http.Request, createdBy string) (*Invite, error) {
	inviteID := id.GetID("cs")
	sURL := &Invite{
		CreatedBy:  createdBy,
		InviteID:   inviteID,
		ExpiryTime: time.Now().Add(s.defaultExpiry),
	}
	err := s.urlDao.CreateInvite(sURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create shortened urls")
	}
	return sURL, nil
}

func (s *impl) GetInviteIfIDisValid(id string) (bool, *Invite) {
	shrtnd := s.urlDao.GetInviteForID(id)
	if len(shrtnd) > 0 && !shrtnd[0].IsExpired() {
		return true, shrtnd[0]
	}
	return false, nil
}
