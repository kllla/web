package invite

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kllla/web/src/config"
	"google.golang.org/api/iterator"
	"log"
)

const bucket = "invites"

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
func (dao *Dao) GetInvites() []*Invite {
	posts := make([]*Invite, 0)
	ctx, cancel := context.WithCancel(dao.ctx)
	iter := dao.client.Collection(bucket).OrderBy("Date", firestore.Desc).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
			cancel()
		}
		post := docToInvite(err, doc, cancel)
		posts = append(posts, post)
	}
	return posts
}

//docToInvite Json marshals firestore doc to Post struct
func docToInvite(err error, doc *firestore.DocumentSnapshot, cancel context.CancelFunc) *Invite {
	md, err := json.Marshal(doc.Data())
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
		cancel()
	}
	var post = &Invite{}
	json.Unmarshal(md, post)
	return post
}

func (dao *Dao) CreateInvite(sURL *Invite) error {
	_, _, err := dao.client.Collection(bucket).Add(dao.ctx, sURL)
	if err != nil {
		return err
	}
	return nil
}

func (dao *Dao) Close() error {
	return dao.client.Close()
}

func (dao *Dao) GetInvitesCreatedBy(createdBy string) []*Invite {
	Invites := make([]*Invite, 0)
	ctx, cancel := context.WithCancel(dao.ctx)
	iter := dao.client.Collection(bucket).Where("CreatedBy", "==", createdBy).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
			cancel()
		}
		invite := docToInvite(err, doc, cancel)
		Invites = append(Invites, invite)
	}
	return Invites
}

func (dao *Dao) GetInviteForID(inviteID string) []*Invite {
	posts := make([]*Invite, 0)
	ctx, cancel := context.WithCancel(dao.ctx)
	iter := dao.client.Collection(bucket).Where("InviteID", "==", inviteID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
			cancel()
		}
		post := docToInvite(err, doc, cancel)
		posts = append(posts, post)
	}
	return posts
}

func (dao *Dao) DeleteInviteByID(id string) error {
	// GetCredentialsForUsername gets NoPassCredentials from the credentials bucket
	ctx, cancel := context.WithCancel(dao.ctx)
	iter := dao.client.Collection(bucket).Where("ID", "==", id).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
				cancel()
			}
		}
		_, err = doc.Ref.Delete(ctx)
		return err
	}
	return fmt.Errorf("no invite found for username %s", id)
}
