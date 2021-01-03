package config

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"os"
	"strconv"
)

const (

)

var (
	Options       = credOption()
	DefaultConfig = defaultConfig()
	TestConfig    = testConfig()
)

func credOption()option.ClientOption {
	credLocation := os.Getenv("CREDS_LOCATION")

	return option.WithCredentialsFile(credLocation)
}

func testConfig() Config {
	return &TestConf{}
}

func defaultConfig() Config {
	return &Conf{Options: Options}
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
	/*client, err := firestore.NewClient(ctx, "test")
	if err != nil {
		log.Fatalf("firebase.NewClient err: %v", err)
	}
	return client*/
	return nil
}

func (config *Conf) ClientFromConfig(ctx context.Context) *firestore.Client {
	projectID := os.Getenv("BUCKET_NAME")
	useCreds := os.Getenv("USE_CREDS")
	fmt.Printf("Bucket Name %s : %t", projectID, useCreds )
	ucred := false
	if useCreds != "" {
		ucred, _ = strconv.ParseBool(useCreds)
	}
	var app = &firebase.App{}
	var err error
	if ucred {
		app, err = firebase.NewApp(ctx, &firebase.Config{
			AuthOverride:     nil,
			DatabaseURL:      "",
			ProjectID:        "infra-prime",
			ServiceAccountID: "108708042762631808425",
			StorageBucket:    "infra-prime.appspot.com",
		})
	} else {
		app, err = firebase.NewApp(ctx, &firebase.Config{
			ProjectID: projectID,
		}, config.Options)
	}
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}
