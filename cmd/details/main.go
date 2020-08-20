package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/andrewkav/viber"
	"github.com/gopher-go/poll_bot"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

//DATASTORE_PROJECT_ID=freeelections2020  DATASTORE_USERS_TABLE=users_dev VIBER_KEY=4bdcdb0d47e7d3db-9e80cfbcec46b16e  ./details

func MustGetDatastoreClient() *datastore.Client {
	data, err := ioutil.ReadFile(os.Getenv("APP_DEFAULT_JSON_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	creds, err := google.CredentialsFromJSON(context.Background(), data, compute.ComputeScope)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(creds.ProjectID)
	DSClient, err := datastore.NewClient(context.Background(), creds.ProjectID)
	if err != nil {
		log.Fatal(err)
	}

	return DSClient
}

func main() {
	/*
		viberKey := os.Getenv("VIBER_KEY")
		fmt.Println(viberKey)
		v := viber.New(viberKey, "Народный опрос", "https://storage.googleapis.com/freeelections2020-img/bot-logo.jpg")
	*/
	//all := []string{}
	var ud *poll_bot.DatastoreUserDAO
	if os.Getenv("DATASTORE_USERS_TABLE") != "" {
		log.Printf("creating datatore user dao, entity kind = %s\n", os.Getenv("DATASTORE_USERS_TABLE"))
		ud = poll_bot.NewDatastoreUserDAO(MustGetDatastoreClient(), os.Getenv("DATASTORE_USERS_TABLE"))
	}

	var users []poll_bot.StorageUser
	_, err := ud.DSClient.GetAll(context.Background(), datastore.NewQuery(ud.EntityKind), &users)
	if err != nil {
		log.Fatalf(err.Error())
	}

	c := 0
	for k, i := range users {
		if c == 25 {
			//do(all, v)
			//all = []string{}
			c = 0
			fmt.Println(k, "Progress")
			//os.Exit(1)
		}
		/*
			det, err := getDetails(i.ID, v)
			if err != nil {
				fmt.Println("Err getting details", err)
				continue
			}
		*/
		level := i.Level
		newLvl := level
		if level >= 6 {
			newLvl = 5
		} else if level == 5 || level == 4 {
			newLvl = level - 1
		}
		err = ud.UpdateUserLvl(&i, newLvl)
		fmt.Println("IDS", i.ID, level, newLvl)
		if err != nil {
			fmt.Println("Error Updating database : ", err)
		}
		fmt.Println(i)
		//all = append(all, i[0])
		c++
	}
	//do(all, v)
	fmt.Println("Done")
}

func getDetails(userId string, v *viber.Viber) (viber.UserDetails, error) {
	details, err := v.UserDetails(userId)
	if err != nil {
		fmt.Println(err)
		return details, err
	}
	return details, nil
}

func do(all []string, v *viber.Viber) {

	//fmt.Println(len(all))
	/*
		   	_, err := v.SendBroadcastTextMessage(all, fmt.Sprintf(`(meds)Друзья, нам нужна ваша помощь!
		   (corn)Наверняка у каждого есть родственники в глубинке, у которых есть смартфон, но они по разным причинам не хотят участвовать в опросах.
		   (diamond)Но сейчас их ответы важны, как никогда!
		   Чтобы составить достоверную картину мнений, нам не хватает именно их!
		   (speaker)Пожалуйста, попросите их поучаствовать в нашем опросе!
		   (lock)Напомним, это абсолютно анонимно и безопасно!
		   (up_graph)Если каждый из вас привлечет хотя бы 3 человека, то нас станет в 4 раза больше и мы сможем составить репрезентативную выборку населения Беларуси.
		   Никто не будет осуждать их за их мнение, ведь наша цель - не переубедить, а показать реальные голоса!
		   (shrug)Объясните им эти простые принципы и помогите разобраться, если им трудно освоить Viber.
		   (prayer_hands)Спасибо!
		   С уважением, команда Народного Опроса.(heart)`))
		if err != nil {
			fmt.Println(err)
		}
	*/
	//fmt.Println(t)

}
