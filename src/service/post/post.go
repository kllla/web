package post

import "time"

type Post struct {
	ID      string
	Author  string
	Title   string
	Content string
	Public  bool
	Date    time.Time
}

func RenderWrap(posts []*Post) []interface{} {
	var intUrls []interface{}
	intUrls = append(intUrls, posts)
	return intUrls
}
