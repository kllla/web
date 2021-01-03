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
	projectID := os.Getenv("PROJECT_ID")
	useCreds := os.Getenv("USE_CREDS")
	fmt.Printf("Bucket Name %s : %t", projectID, useCreds )
	ucred := false
	if useCreds != "" {
		ucred, _ = strconv.ParseBool(useCreds)
	}
	var app = &firebase.App{}
	var err error
	if ucred {
		type Config struct {
			AuthOverride     *map[string]interface{} `json:"databaseAuthVariableOverride"`
			DatabaseURL      string                  `json:"databaseURL"`
			ProjectID        string                  `json:"projectId"`
			ServiceAccountID string                  `json:"serviceAccountId"`
			StorageBucket    string                  `json:"storageBucket"`
		}
		app, err = firebase.NewApp(ctx, &firebase.Config{
			AuthOverride:     nil,
			DatabaseURL:      "",
			ProjectID:        "infra-person",
			ServiceAccountID: "105637127689182478722",
			StorageBucket:    "bucket",
		}, config.Options)
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
