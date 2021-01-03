package shorten

import "time"

type ShortenedURL struct {
	CreatedBy    string
	ShortenedID  string
	ShortenedURL string
	LongUrl      string
	ExpiryTime   time.Time
}

func (s *ShortenedURL) IsExpired() bool {
	return time.Now().After(s.ExpiryTime)
}

func RenderWrap(shortenedURLS []*ShortenedURL) []interface{} {
	var intUrls []interface{}
	intUrls = append(intUrls, shortenedURLS)
	return intUrls
}
