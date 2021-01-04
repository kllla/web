package credentials

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kllla/web/src/config"
	"google.golang.org/api/iterator"
	"log"
)

const credentialsBucket = "credentials"

type Dao interface {
	SetClientFromConfig(config config.Config)
	CreateCredentials(credentials *Credentials) error
	DeleteCredentials(credentials *Credentials) error
	GetAllCredentials() []*NoPassCredentials
	GetCredentialsForUsername(username string) []*NoPassCredentials
	Close() error
}

type daoImpl struct {
	client *firestore.Client
	ctx    context.Context
}

func (dao *daoImpl) SetClientFromConfig(config config.Config) {
	dao.client = config.ClientFromConfig(context.Background())
}

func (dao *daoImpl) DeleteCredentials(credentials *Credentials) error {
	return dao.DeleteCredentialsForUsername(credentials.Username)
}

func NewDao(config config.Config) Dao {
	ctx := context.Background()
	d := &daoImpl{
		ctx: ctx,
	}
	d.SetClientFromConfig(config)
	return d
}

func (dao *daoImpl) CreateCredentials(credentials *Credentials) error {
	existingCredentials := dao.GetCredentialsForUsername(credentials.Username)
	if len(existingCredentials) > 0 {
		return fmt.Errorf("username %s is unavailable", credentials.Username)
	}
	_, _, err := dao.client.Collection(credentialsBucket).Add(dao.ctx, credentials.ToNoPassCredentials())
	if err != nil {
		log.Fatalf("Failed creating credentials: %s", err)
	}
	return nil
}

// GetCredentialsForUsername gets NoPassCredentials from the credentials bucket
func (dao *daoImpl) DeleteCredentialsForUsername(username string) error {
	iter := dao.client.Collection(credentialsBucket).Where("Username", "==", username).Documents(dao.ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			failToIterate(err)
		}
		_, err = doc.Ref.Delete(dao.ctx)
		return err
	}
	return fmt.Errorf("no credentials found for username %s", username)
}

// GetCredentialsForUsername gets NoPassCredentials from the credentials bucket
func (dao *daoImpl) GetCredentialsForUsername(username string) []*NoPassCredentials {
	credentials := make([]*NoPassCredentials, 0)
	iter := dao.client.Collection(credentialsBucket).Where("Username", "==", username).Documents(dao.ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			failToIterate(err)
		}
		cred := docToCredentials(err, doc)
		credentials = append(credentials, cred)
	}
	return credentials
}

func failToIterate(err error) {
	if err != nil {
		log.Fatalf("Failed to iterate: %v", err)
	}
}

// GetAllCredentials gets all posts from
func (dao *daoImpl) GetAllCredentials() []*NoPassCredentials {
	credentials := make([]*NoPassCredentials, 0)
	iter := dao.client.Collection(credentialsBucket).Documents(dao.ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			failToIterate(err)
		}
		post := docToCredentials(err, doc)
		credentials = append(credentials, post)
	}
	return credentials
}

//docToCredentials Json marshals firestore doc to Post struct
func docToCredentials(err error, doc *firestore.DocumentSnapshot) *NoPassCredentials {
	md, err := json.Marshal(doc.Data())
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}
	var credentials = &NoPassCredentials{}
	json.Unmarshal(md, credentials)
	return credentials
}

func (dao *daoImpl) Close() error {
	return dao.client.Close()
}
