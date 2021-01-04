package sao

import (
	"cloud.google.com/go/storage"
	"context"
	"io/ioutil"
	"log"
	"os"
)

// Sao is the default actions required for all Saos
// independent of their bucket
type Sao interface {
	GetStaticFiles(object string) []byte
}

type saoImpl struct {
	ctx    context.Context
	client *storage.Client
}

// New returns the impl for the Sao interface after
// initialising the context and client
func New() Sao {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return &saoImpl{ctx: ctx, client: client}
}

// GetStaticFiles is the primary function to retrieve file
// object from the bucket
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

func (sao *saoImpl) close() error {
	return sao.client.Close()
}
