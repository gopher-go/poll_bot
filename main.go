package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mileusna/viber"
	"log"
	"os"
)

var caseSensitive = true

func main() {
	err := execute()
	if err != nil {
		log.Fatal(err)
	}
	<-make(chan int)
}

func execute() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Problem with .env file: %v", err)
	}

	viberKey := os.Getenv("VIBER_KEY")
	callback_URL := os.Getenv("CALLBACK_URL")

	var ud userDAO
	ud, err = newPQUserDAO(os.Getenv("DB_CONNECTION"))
	if err != nil {
		log.Fatal(err)
	}

	v := viber.New(viberKey, "Voting bot", "https://thumbs.dreamstime.com/z/human-hand-write-yes-vote-voting-paper-pen-flat-concept-illustration-man-s-red-pen-ballot-check-sign-88802664.jpg")
	go func() {
		err := serve(v, ud)
		if err != nil {
			log.Fatal(err)
		}
	}()
	_, err = v.SetWebhook(callback_URL, nil)
	return err
}
