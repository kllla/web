package shorten

import (
	"fmt"
	"github.com/kllla/web/src/config"
	"github.com/kllla/web/src/service/id"
	"net/http"
	"time"
)

const domain = "kll.la/"

type Service interface {
	CreatedShortenedURL(w http.ResponseWriter, r *http.Request, createdBy string) (*ShortenedURL, error)
	GetUrlIfIDisValid(id string) (bool, *ShortenedURL)
	GetAllShortenedURLsCreatedBy(by string) []*ShortenedURL
}

type impl struct {
	urlDao        *Dao
	defaultExpiry time.Duration
}

func (s *impl) GetURLFromFormData(w http.ResponseWriter, r *http.Request) string {
	err := r.ParseForm()
	if err != nil {
		return ""
	}
	longURL := r.FormValue("longURL")

	return longURL
}

func (s *impl) GetAllShortenedURLsCreatedBy(createdBy string) []*ShortenedURL {
	return s.urlDao.GetShortenedURLsCreatedBy(createdBy)
}

const aWeek = (time.Hour * 24) * 7

func NewService() Service {
	return &impl{
		urlDao:        NewDao(config.DefaultConfig),
		defaultExpiry: aWeek,
	}
}

func (s *impl) CreatedShortenedURL(w http.ResponseWriter, r *http.Request, createdBy string) (*ShortenedURL, error) {
	url := s.GetURLFromFormData(w, r)
	shrtnID := id.GetID(url)
	sURL := &ShortenedURL{
		CreatedBy:    createdBy,
		ShortenedID:  shrtnID,
		ShortenedURL: domain + shrtnID,
		LongUrl:      url,
		ExpiryTime:   time.Now().Add(s.defaultExpiry),
	}
	err := s.urlDao.CreateShortenedURL(sURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create shortened urls")
	}
	return sURL, nil
}

func (s *impl) GetUrlIfIDisValid(id string) (bool, *ShortenedURL) {
	shrtnd := s.urlDao.GetShortenedURLForID(id)
	if len(shrtnd) > 0 && !shrtnd[0].IsExpired() {
		return true, shrtnd[0]
	}
	return false, nil
}
