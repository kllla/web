package post

import (
	"fmt"
	"github.com/kllla/web/src/config"
	"github.com/kllla/web/src/service/id"
	"net/http"
	"strconv"
	"time"
)

type Service interface {
	CreatePost(post *Post)
	DeletePost(id string) bool
	GetAllPostsForUsername(username string) []*Post
	GetPostByID(pathID string) []*Post
	GetPosts() []*Post
	GetPublicPosts() []*Post
	GetHiddenPosts() []*Post
	GetActionFromFormData(w http.ResponseWriter, r *http.Request) string
	GetPostFromFormData(w http.ResponseWriter, r *http.Request, author string) *Post
}

type impl struct {
	postDao Dao
}

func NewService() Service {
	return &impl{postDao: NewDao(config.DefaultConfig)}
}

func (m *impl) GetPosts() []*Post {
	return m.postDao.GetPosts()
}

func (m *impl) GetPublicPosts() []*Post {
	return m.postDao.GetPublicPosts()
}

func (m *impl) GetHiddenPosts() []*Post {
	return m.postDao.GetHiddenPosts()
}

func (m *impl) CreatePost(post *Post) {
	m.postDao.CreatePost(post)
}
func (m *impl) GetActionFromFormData(w http.ResponseWriter, r *http.Request) string {
	if r.ParseForm() != nil {
		return ""
	}
	action := r.FormValue("action")
	return action
}

func (m *impl) GetPostFromFormData(w http.ResponseWriter, r *http.Request, author string) *Post {
	if r.ParseForm() != nil {
		return nil
	}
	title := r.FormValue("title")
	content := r.FormValue("content")
	publicForm := r.FormValue("public")
	public, err := strconv.ParseBool(publicForm)
	if err != nil {
		fmt.Print("failed publiccheck")
		public = false
	}
	return &Post{
		ID:      id.GetID(title),
		Author:  author,
		Title:   title,
		Content: content,
		Public:  public,
		Date:    time.Now(),
	}
}

func (m *impl) GetAllPostsForUsername(username string) []*Post {
	return m.postDao.GetPostsForUsername(username)
}

func (m *impl) DeletePost(id string) bool {
	return m.postDao.DeletePostById(id)
}

func (m *impl) GetPostByID(pathID string) []*Post {
	return m.postDao.GetPostByID(pathID)
}
