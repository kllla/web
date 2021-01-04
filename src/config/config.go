package config

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"os"
)

var (
	Options       = credOption()
	DefaultConfig = defaultConfig()
	TestConfig    = testConfig()
)

func credOption() option.ClientOption {
	credLocation := os.Getenv("CREDS_LOCATION")
	log.Println("CRED LOCATION", credLocation)
	w, _ := os.Getwd()
	log.Println("PWD = ", w)
	return option.WithCredentialsFile(credLocation)
}

func testConfig() Config {
	return &TestConf{}
}

func defaultConfig() Config {
	return &Conf{Options: nil}
}

type Config interface {
	ClientFromConfig(ctx context.Context) *firestore.Client
	GetOptions() option.ClientOption
}

type Conf struct {
	Options option.ClientOption
}

func (config *Conf) GetOptions() option.ClientOption {
	return config.Options
}

type TestConf struct {
}

func (conf *TestConf) GetOptions() option.ClientOption {
	return nil
}

func (conf *TestConf) ClientFromConfig(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, "test")
	if err != nil {
		log.Fatalf("firebase.NewClient err: %v", err)
	}
	return client
}

func (config *Conf) ClientFromConfig(ctx context.Context) *firestore.Client {
	projectID := os.Getenv("PROJECT_ID")
	serviceAccountID := os.Getenv("SERVICE_ACCOUNT_ID")
	storageBucket := os.Getenv("BUCKET_NAME")
	app, err := firebase.NewApp(ctx, &firebase.Config{
		AuthOverride:     nil,
		DatabaseURL:      "",
		ProjectID:        projectID,
		ServiceAccountID: serviceAccountID,
		StorageBucket:    storageBucket,
	})
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}
