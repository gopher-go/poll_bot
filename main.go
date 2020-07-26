package main

import (
	"fmt"
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
	v := viber.New(viberKey, "Народный опрос", "https://thumbs.dreamstime.com/z/human-hand-write-yes-vote-voting-paper-pen-flat-concept-illustration-man-s-red-pen-ballot-check-sign-88802664.jpg")
	go func() {
		err := serve(v, ud)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return setViberWebhook(v, callbackURL)
}
