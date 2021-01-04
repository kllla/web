package sao

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"os"
)

type Sao interface {
	GetStaticFiles(object string) []byte
	Close() error
}

type saoImpl struct {
	ctx    context.Context
	client *storage.Client
}

type Config struct {
	CredentialsPath string
}

func NewSao() Sao {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("CREDS_LOCATION")))
	if err != nil {
		log.Fatalln(err)
	}
	return &saoImpl{ctx: ctx, client: client}
}

func (sao *saoImpl) GetStaticFiles(object string) []byte {
	var bucket = os.Getenv("BUCKET_NAME")
	rc, err := sao.client.Bucket(bucket).Object(object).NewReader(sao.ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer rc.Close()
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Fatalln(err)
	}
	return data

}

func (sao *saoImpl) Close() error {
	return sao.client.Close()
}
