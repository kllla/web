package post

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"github.com/kllla/web/src/config"
	"google.golang.org/api/iterator"
	"log"
)

const bucket = "posts"

type Dao interface {
	SetClientFromConfig(config config.Config)
	CreatePost(post *Post)
	GetPosts() []*Post
	GetHiddenPosts() []*Post
	GetPublicPosts() []*Post
	GetPostsForUsername(username string) []*Post
	Close() error
	DeletePostById(id string) bool
	GetPostByID(id string) []*Post
}

type daoImpl struct {
	client *firestore.Client
	ctx    context.Context
}

func (dao *daoImpl) GetPostByID(id string) []*Post {
	posts := make([]*Post, 0)
	iter := dao.client.Collection(bucket).Where("ID", "==", id).Documents(dao.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		post := docToPost(err, doc)
		posts = append(posts, post)
	}
	return posts
}

func NewDao(config config.Config) Dao {
	ctx := context.Background()
	client := config.ClientFromConfig(ctx)
	d := &daoImpl{
		client: client,
		ctx:    ctx,
	}
	d.SetClientFromConfig(config)
	return d
}

// GetPosts gets all posts
func (dao *daoImpl) GetPosts() []*Post {
	posts := make([]*Post, 0)
	iter := dao.client.Collection(bucket).OrderBy("Date", firestore.Desc).Documents(dao.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		post := docToPost(err, doc)
		posts = append(posts, post)
	}
	return posts
}

func (dao *daoImpl) DeletePostById(id string) bool {
	iter := dao.client.Collection(bucket).Where("ID", "==", id).Documents(dao.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		_, err = doc.Ref.Delete(dao.ctx)
		if err != nil {
			log.Fatalf("Failed to delete: %v", err)
			return false
		}
	}
	return true
}

// GetPostsForUsername gets all posts for the username
func (dao *daoImpl) GetPostsForUsername(username string) []*Post {
	posts := make([]*Post, 0)
	iter := dao.client.Collection(bucket).Where("Author", "==", username).OrderBy("Date", firestore.Desc).Documents(dao.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		post := docToPost(err, doc)
		posts = append(posts, post)
	}
	return posts
}

// GetHiddenPosts gets all hidden posts
func (dao *daoImpl) GetHiddenPosts() []*Post {
	posts := make([]*Post, 0)
	iter := dao.client.Collection(bucket).Where("Public", "==", "false").OrderBy("Date", firestore.Desc).Documents(dao.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		post := docToPost(err, doc)
		posts = append(posts, post)
	}
	return posts
}

// GetPublicPosts gets all posts marked public
func (dao *daoImpl) GetPublicPosts() []*Post {
	posts := make([]*Post, 0)
	iter := dao.client.Collection(bucket).Where("Public", "==", true).OrderBy("Date", firestore.Desc).Documents(dao.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		post := docToPost(err, doc)
		posts = append(posts, post)
	}
	return posts
}

func (dao *daoImpl) CreatePost(post *Post) {
	_, _, err := dao.client.Collection(bucket).Add(dao.ctx, post)
	if err != nil {
		log.Fatalf("Failed adding posting: %v", err)
	}
}

func (dao *daoImpl) SetClientFromConfig(config config.Config) {
	dao.client = config.ClientFromConfig(context.Background())
}

func (dao *daoImpl) Close() error {
	return dao.client.Close()
}

//docToPost Json marshals firestore doc to Post struct
func docToPost(err error, doc *firestore.DocumentSnapshot) *Post {
	md, err := json.Marshal(doc.Data())
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}
	var post = &Post{}
	json.Unmarshal(md, post)
	return post
}
