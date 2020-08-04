package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/andrewkav/viber"
	"github.com/gopher-go/poll_bot"
)

func main() {
	viberKey := os.Getenv("VIBER_KEY")
	fmt.Println(viberKey)
	v := viber.New(viberKey, "Народный опрос", "https://storage.googleapis.com/freeelections2020-img/bot-logo.jpg")
	//all := []string{}
	var ud *poll_bot.DatastoreUserDAO
	if os.Getenv("DATASTORE_USERS_TABLE") != "" {
		log.Printf("creating datatore user dao, entity kind = %s\n", os.Getenv("DATASTORE_USERS_TABLE"))
		ud = poll_bot.NewDatastoreUserDAO(poll_bot.MustGetDatastoreClient(), os.Getenv("DATASTORE_USERS_TABLE"))
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
			os.Exit(1)
		}
		det := getDetails(i.ID, v)
		ud.Update(&i, det.Mcc, det.Mnc)
		fmt.Println(i)
		//all = append(all, i[0])
		c++
	}
	//do(all, v)
	fmt.Println("Done")
}

func getDetails(userId string, v *viber.Viber) viber.UserDetails {
	details, err := v.UserDetails(userId)
	if err != nil {
		fmt.Println(err)
		return details
	}
	return details
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
