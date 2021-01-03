package invite

import "time"

type Invite struct {
	CreatedBy  string
	InviteID   string
	ExpiryTime time.Time
}

func (s *Invite) IsExpired() bool {
	return time.Now().After(s.ExpiryTime)
}

func RenderWrap(invites []*Invite) []interface{} {
	var intUrls []interface{}
	intUrls = append(intUrls, invites)
	return intUrls
}
