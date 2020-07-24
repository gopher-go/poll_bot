package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/mileusna/viber"
)

var caseSensitive = true

func main() {
	err := execute()
	if err != nil {
		log.Fatal(err)
	}
	<-make(chan int)
}

func init() {
	newUserDAO = func() (userDAO, error) {
		ud, err := newPQUserDAO()
		return ud, err
	}
}

func execute() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Problem with .env file: %v", err)
	}

	viberKey := os.Getenv("VIBER_KEY")
	callback_URL := os.Getenv("CALLBACK_URL")

	// setup persistence
	// prefer redis persistence over SQL
	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		var (
			rudOnce sync.Once
			rud     *redisUserDAO
		)
		newUserDAO = func() (userDAO, error) {
			rudOnce.Do(func() {
				rud = newRedisUserDAO(redis.NewClient(&redis.Options{Addr: redisAddr}))
			})

			return rud, nil
		}
	}

	v := viber.New(viberKey, "Voting bot", "https://thumbs.dreamstime.com/z/human-hand-write-yes-vote-voting-paper-pen-flat-concept-illustration-man-s-red-pen-ballot-check-sign-88802664.jpg")
	go func() {
		err := serve(v)
		if err != nil {
			log.Fatal(err)
		}
	}()
	_, err = v.SetWebhook(callback_URL, nil)
	return err
}
