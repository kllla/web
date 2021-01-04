package shorten

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"github.com/kllla/web/src/config"
	"google.golang.org/api/iterator"
	"log"
)

const bucket = "shorten"

type Dao struct {
	client *firestore.Client
	ctx    context.Context
}

func NewDao(config config.Config) *Dao {
	ctx := context.Background()
	client := config.ClientFromConfig(ctx)
	return &Dao{
		client: client,
		ctx:    ctx,
	}
}

// getPosts gets all posts from
func (dao *Dao) GetShortenedURLS() []*ShortenedURL {
	posts := make([]*ShortenedURL, 0)
	iter := dao.client.Collection(bucket).OrderBy("Date", firestore.Desc).Documents(dao.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		post := docToShortenedURL(err, doc)
		posts = append(posts, post)
	}
	return posts
}

//docToShortenedURL Json marshals firestore doc to Post struct
func docToShortenedURL(err error, doc *firestore.DocumentSnapshot) *ShortenedURL {
	md, err := json.Marshal(doc.Data())
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}
	var post = &ShortenedURL{}
	json.Unmarshal(md, post)
	return post
}

func (dao *Dao) CreateShortenedURL(sURL *ShortenedURL) error {
	_, _, err := dao.client.Collection(bucket).Add(dao.ctx, sURL)
	if err != nil {
		return err
	}
	return nil
}

func (dao *Dao) Close() error {
	return dao.client.Close()
}

func (dao *Dao) GetShortenedURLsCreatedBy(createdBy string) []*ShortenedURL {
	ShortenedURLS := make([]*ShortenedURL, 0)
	iter := dao.client.Collection(bucket).Where("CreatedBy", "==", createdBy).Documents(dao.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		ShortenedURL := docToShortenedURL(err, doc)
		ShortenedURLS = append(ShortenedURLS, ShortenedURL)
	}
	return ShortenedURLS
}

func (dao *Dao) GetShortenedURLForID(shortenedID string) []*ShortenedURL {
	posts := make([]*ShortenedURL, 0)
	iter := dao.client.Collection(bucket).Where("ShortenedID", "==", shortenedID).Documents(dao.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		post := docToShortenedURL(err, doc)
		posts = append(posts, post)
	}
	return posts
}
