package poll_bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/olivere/elastic/v7"

	"cloud.google.com/go/datastore"
	"github.com/andrewkav/viber"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

var caseSensitive = true

func main() {
	err := execute()
	if err != nil {
		log.Fatal(err)
	}
	<-make(chan int)
}

func setViberWebhook(v *viber.Viber, url string) error {
	req := viber.WebhookReq{
		URL:        url,
		EventTypes: nil,
		SendName:   false,
		SendPhoto:  false,
	}
	_, err := v.PostData("https://chatapi.viber.com/pa/set_webhook", req)
	if err != nil {
		return err
	}

	return nil
}

func MustGetDatastoreClient() *datastore.Client {
	creds, err := google.FindDefaultCredentials(context.Background(), compute.ComputeScope)
	if err != nil {
		log.Fatal(err)
	}
	DSClient, err := datastore.NewClient(context.Background(), creds.ProjectID)
	if err != nil {
		log.Fatal(err)
	}

	return DSClient
}

func execute() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Problem with .env file: %v", err)
	}

	viberKey := os.Getenv("VIBER_KEY")
	callbackURL := os.Getenv("CALLBACK_URL")

	var ud userDAO
	if os.Getenv("DATASTORE_USERS_TABLE") != "" {
		log.Printf("creating datatore user dao, entity kind = %s\n", os.Getenv("DATASTORE_USERS_TABLE"))
		ud = NewDatastoreUserDAO(MustGetDatastoreClient(), os.Getenv("DATASTORE_USERS_TABLE"))
	} else {
		ud, err = newPQUserDAO(os.Getenv("DB_CONNECTION"))
		if err != nil {
			log.Fatal(err)
		}
	}

	var ld logDAO
	if os.Getenv("DATASTORE_USER_ANSWER_LOG_TABLE") != "" {
		log.Printf("creating datastore log answer dao, entity kind = %s\n", os.Getenv("DATASTORE_USER_ANSWER_LOG_TABLE"))
		ld = newDatastoreLogDAO(MustGetDatastoreClient(), os.Getenv("DATASTORE_USER_ANSWER_LOG_TABLE"))
	}

	var sd *statsDao
	if os.Getenv("ELASTIC_HOSTS") != "" {
		log.Printf("creating stats dao, ES_HOST=%s\n", os.Getenv("DATASTORE_USER_ANSWER_LOG_TABLE"))
		esClient, err := elastic.NewSimpleClient(elastic.SetURL(strings.Split(os.Getenv("ELASTIC_HOSTS"), ",")...),
			elastic.SetBasicAuth(os.Getenv("ELASTIC_BASIC_AUTH_USER"), os.Getenv("ELASTIC_BASIC_AUTH_PASSWORD")))
		if err != nil {
			log.Printf("unable to create ES client, err=%v\n", err)
		} else {
			sd = newStatsDao(esClient)
		}
	}

	v := viber.New(viberKey, "Народный опрос", "https://storage.googleapis.com/freeelections2020-img/bot-logo.jpg")
	go func() {
		err := serve(v, ud, ld, sd)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return setViberWebhook(v, callbackURL)
}
