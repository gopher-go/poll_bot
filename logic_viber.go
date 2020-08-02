package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

var (
	errPleaseChooseSuggestedAnswer = errors.New("Пожалуйста, выберите один из предложенных вариантов ответа.")
)

func knownEvent(c *ViberCallback) bool {
	return c.Event == "message" ||
		c.Event == "delivered" ||
		c.Event == "seen" ||
		c.Event == "subscribed" ||
		c.Event == "unsubscribed" ||
		c.Event == "conversation_started" ||
		c.Event == "webhook"
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type viberReply struct {
	text    string
	options []string
}

func logUserAnswer(s *storage, al *answerLog) {
	go func() {
		if err := s.LogAnswer(al); err != nil {
			log.Printf("unable to log user answer, err=%v\n", err)
		}
	}()
}

func generateReplyFor(poll poll, s *storage, callback *ViberCallback) (*viberReply, error) {
	if !knownEvent(callback) {
		return nil, fmt.Errorf("Unknown message %v", callback.Event)
	}

	if strings.ToLower(callback.Message.Text) == "i'm a tester, clear it" {
		err := s.Clear(callback.User.ID)
		return &viberReply{text: fmt.Sprintf("Your storage cleared with %v", err)}, nil
	}

	storageUser, err := s.Obtain(callback.User.ID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if storageUser.isChanged {
			_ = s.Persist(storageUser)
		}
	}()

	if storageUser.Country == "" && callback.User.Country != "" {
		storageUser.Country = callback.User.Country
		storageUser.isChanged = true
	}
	if storageUser.Language == "" && callback.User.Language != "" {
		storageUser.Language = callback.User.Language
		storageUser.isChanged = true
	}
	if storageUser.MobileNetworkCode == 0 && callback.User.MNC != 0 {
		storageUser.MobileNetworkCode = callback.User.MNC
		storageUser.isChanged = true
	}
	if storageUser.MobileCountryCode == 0 && callback.User.MCC != 0 {
		storageUser.MobileCountryCode = callback.User.MCC
		storageUser.isChanged = true
	}
	if storageUser.Context == "" && callback.Context != "" {
		storageUser.Context = callback.Context
		storageUser.isChanged = true
	}

	if callback.Event == "unsubscribed" {
		storageUser.Properties["ConversationStarted"] = "false"
		storageUser.isChanged = true
		return nil, nil
	}

	if callback.Event == "message" {
		pi := poll.getLevel(storageUser.Level)

		al := answerLog{
			UserID:      storageUser.ID,
			UserContext: storageUser.Context,
			QuestionID:  pi.id,
			Answer:      callback.Message.Text,
			AnswerLevel: pi.level,
			IsValid:     true,
			CreatedAt:   time.Now().UTC(),
		}

		err := analyseAnswer(pi, storageUser, callback)
		if err != nil {
			reply, _ := getViberReplyForLevel(poll, s, storageUser, callback)

			// if finished don't generate error
			if !poll.isFinishedFor(storageUser) {
				reply.text = err.Error() + "\n\n" + reply.text
			}

			al.IsValid = false
			logUserAnswer(s, &al)

			return reply, nil
		}

		storageUser.Level++
		storageUser.isChanged = true

		logUserAnswer(s, &al)

		return getViberReplyForLevel(poll, s, storageUser, callback)
	}

	if storageUser.Properties["ConversationStarted"] != "true" {
		reply, err := getViberReplyForLevel(poll, s, storageUser, callback)
		if err != nil {
			return nil, err
		}
		storageUser.Properties["ConversationStarted"] = "true"
		storageUser.isChanged = true
		return reply, nil
	}

	return nil, nil
}

const url = "narodny-opros.info"
const welcomeHeader = `Добро пожаловать в проект «Народный опрос»! 

Давайте вместе узнаем реальный предвыборный рейтинг всех кандидатов в президенты!
Всё, что необходимо сделать, — ответить на несколько вопросов.
Не беспокойтесь, ваше участие полностью анонимное.

Нас уже %d человек! Присоединяйтесь!

Чтобы ответить на вопрос, выберите один из предложенных вариантов.`

func getViberReplyForLevel(p poll, s *storage, u *storageUser, c *ViberCallback) (*viberReply, error) {

	var isNewConversation = u.Properties["ConversationStarted"] != "true"

	if p.isFinishedFor(u) {
		totalCount, err := s.CountCached()
		if err != nil {
			return nil, err
		}
		text := "Спасибо за участие в нашем опросе!\nСледите за динамикой опроса на сайте " + url
		if isNewConversation {
			text = "Добрый день!\nСпасибо за участие в нашем опросе!\nСледите за динамикой опроса на сайте " + url
		}
		text += fmt.Sprintf("\nНас уже %d человек!", totalCount)
		return &viberReply{text: text}, nil
	}

	var welcome string
	if isNewConversation {
		totalCount, err := s.CountCached()
		if err != nil {
			return nil, err
		}
		welcome = fmt.Sprintf(welcomeHeader, totalCount) + "\n\n"
	}

	item := p.getLevel(u.Level)

	return &viberReply{
		text:    welcome + item.question(u, c),
		options: item.possibleAnswers,
	}, nil
}

func analyseAnswer(pi pollItem, u *storageUser, c *ViberCallback) error {
	answer := c.Message.Text
	normalAnswer := answer
	if !caseSensitive {
		found := false
		answer = strings.ToLower(answer)

		// handle click reply
		for _, v := range pi.possibleAnswers {
			if answer == strings.ToLower(v) {
				normalAnswer = v
				found = true
				break
			}
		}

		if !found {
			return errPleaseChooseSuggestedAnswer
		}
	} else if pi.possibleAnswers != nil && !contains(pi.possibleAnswers, answer) {
		return errPleaseChooseSuggestedAnswer
	}

	if pi.validateAnswer != nil {
		err := pi.validateAnswer(normalAnswer)
		if err != nil {
			return err
		}
	}
	if pi.persistAnswer != nil {
		err := pi.persistAnswer(normalAnswer, u)
		if err != nil {
			return err
		}
	}
	return nil
}
