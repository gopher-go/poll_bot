package main

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"log"
	"os"

	"github.com/andrewkav/viber"
	"github.com/joho/godotenv"
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

	return err

}

func execute() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Problem with .env file: %v", err)
	}

	viberKey := os.Getenv("VIBER_KEY")
	callbackURL := os.Getenv("CALLBACK_URL")

	var ud userDAO
	ud, err = newPQUserDAO(os.Getenv("DB_CONNECTION"))
	if err != nil {
		log.Fatal(err)
	}

	var ld logDAO
	if os.Getenv("DATASTORE_USER_ANSWER_LOG_TABLE") != "" {
		creds, err := google.FindDefaultCredentials(context.Background(), compute.ComputeScope)
		if err != nil {
			log.Fatal(err)
		}
		dsClient, err := datastore.NewClient(context.Background(), creds.ProjectID)
		if err != nil {
			log.Fatal(err)
		}

		ld = newDatastoreLogDAO(dsClient, os.Getenv("DATASTORE_USER_ANSWER_LOG_TABLE"))
	}

	v := viber.New(viberKey, "Народный опрос", "https://storage.googleapis.com/freeelections2020-img/bot-logo.jpg")
	go func() {
		err := serve(v, ud, ld)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return setViberWebhook(v, callbackURL)
}
