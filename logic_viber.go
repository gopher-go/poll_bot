package poll_bot

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

	StorageUser, err := s.Obtain(callback.User.ID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if StorageUser.isChanged {
			_ = s.Persist(StorageUser)
		}
	}()

	if StorageUser.Country == "" && callback.User.Country != "" {
		StorageUser.Country = callback.User.Country
		StorageUser.isChanged = true
	}
	if StorageUser.Language == "" && callback.User.Language != "" {
		StorageUser.Language = callback.User.Language
		StorageUser.isChanged = true
	}
	if StorageUser.MobileNetworkCode == 0 && callback.User.MNC != 0 {
		StorageUser.MobileNetworkCode = callback.User.MNC
		StorageUser.isChanged = true
	}
	if StorageUser.MobileCountryCode == 0 && callback.User.MCC != 0 {
		StorageUser.MobileCountryCode = callback.User.MCC
		StorageUser.isChanged = true
	}
	if StorageUser.Context == "" && callback.Context != "" {
		StorageUser.Context = callback.Context
		StorageUser.isChanged = true
	}

	if callback.Event == "unsubscribed" {
		StorageUser.Properties["ConversationStarted"] = "false"
		StorageUser.isChanged = true
		return nil, nil
	}

	if callback.Event == "message" {
		pi := poll.getLevel(StorageUser.Level)

		al := answerLog{
			UserID:      StorageUser.ID,
			UserContext: StorageUser.Context,
			QuestionID:  pi.id,
			Answer:      callback.Message.Text,
			AnswerLevel: pi.level,
			IsValid:     true,
			CreatedAt:   time.Now().UTC(),
		}

		err := analyseAnswer(pi, StorageUser, callback)
		if err != nil {
			reply, _ := getViberReplyForLevel(poll, s, StorageUser, callback)

			// if finished don't generate error
			if !poll.isFinishedFor(StorageUser) {
				reply.text = err.Error() + "\n\n" + reply.text
			}

			al.IsValid = false
			logUserAnswer(s, &al)

			return reply, nil
		}

		StorageUser.Level++
		StorageUser.isChanged = true

		logUserAnswer(s, &al)

		return getViberReplyForLevel(poll, s, StorageUser, callback)
	}

	if StorageUser.Properties["ConversationStarted"] != "true" {
		reply, err := getViberReplyForLevel(poll, s, StorageUser, callback)
		if err != nil {
			return nil, err
		}
		StorageUser.Properties["ConversationStarted"] = "true"
		StorageUser.isChanged = true
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

func getViberReplyForLevel(p poll, s *storage, u *StorageUser, c *ViberCallback) (*viberReply, error) {

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

func analyseAnswer(pi pollItem, u *StorageUser, c *ViberCallback) error {
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
