package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	errPleaseChooseSuggestedAnswer = errors.New("Пожалуйста выберите один из предложенных ответов или введите его номер.")
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

func generateReplyFor(poll poll, storage *storage, callback *ViberCallback) (*viberReply, error) {
	if !knownEvent(callback) {
		return nil, fmt.Errorf("Unknown message %v", callback.Event)
	}

	if strings.ToLower(callback.Message.Text) == "i'm a tester, clear it" {
		err := storage.Clear(callback.User.ID)
		return &viberReply{text: fmt.Sprintf("Your storage cleared with %v", err)}, nil
	}

	storageUser, err := storage.Obtain(callback.User.ID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if storageUser.isChanged {
			_ = storage.Persist(storageUser)
		}
	}()

	if storageUser.Country == "" && callback.User.Country != "" {
		storageUser.Country = callback.User.Country
		storageUser.isChanged = true
	}

	if callback.Event == "unsubscribed" {
		storageUser.Properties["ConversationStarted"] = "false"
		storageUser.isChanged = true
		return nil, nil
	}

	if callback.Event == "conversation_started" {
		storageUser.Context = callback.Context
		storageUser.isChanged = true
	}

	if callback.Event == "message" {
		err := analyseAnswer(poll, storageUser, callback)
		if err != nil {
			reply, _ := getViberReplyForLevel(poll, storage, storageUser, callback)
			reply.text = err.Error() + "\n\n" + reply.text
			return reply, nil
		}
		storageUser.Level++
		storageUser.isChanged = true
		return getViberReplyForLevel(poll, storage, storageUser, callback)
	}

	if storageUser.Properties["ConversationStarted"] != "true" {
		reply, err := getViberReplyForLevel(poll, storage, storageUser, callback)
		if err != nil {
			return nil, err
		}
		storageUser.Properties["ConversationStarted"] = "true"
		storageUser.isChanged = true
		return reply, nil
	}

	return nil, nil
}

func getViberReplyForLevel(p poll, s *storage, u *storageUser, c *ViberCallback) (*viberReply, error) {
	var welcome string
	if u.Properties["ConversationStarted"] != "true" {
		welcome = "Добрый день!\n"
	}

	if p.isFinishedFor(u) {
		totalCount, err := s.PersistCount()
		if err != nil {
			return nil, err
		}
		text := "Спасибо, ваш голос учтен!"
		if welcome != "" {
			text = welcome + "Вы уже приняли участие в Народном опросе. " + text
		}
		if totalCount > 500 {
			text += fmt.Sprintf("\nНас уже %d человек!", totalCount)
		}
		return &viberReply{text: text}, nil
	}

	item := p.getLevel(u.Level)
	reply := viberReply{text: fmt.Sprintf("Непонятно. Нет уровня %v в вопросах", u.Level)}
	if item != nil {
		reply.text = welcome + item.question(u, c)
		reply.options = item.possibleAnswers
	}
	return &reply, nil
}

func analyseAnswer(p poll, u *storageUser, c *ViberCallback) error {
	item := p.getLevel(u.Level)
	if item == nil {
		return nil
	}

	answer := c.Message.Text
	normalAnswer := answer

	found := false
	// handle numeric reply
	if n, err := strconv.Atoi(answer); err == nil {
		if n > len(item.possibleAnswers) || n < 1 {
			return errPleaseChooseSuggestedAnswer
		}
		normalAnswer = item.possibleAnswers[n-1]
		found = true
	}

	if !found {
		if !caseSensitive {
			answer = strings.ToLower(answer)

			// handle click reply
			for _, v := range item.possibleAnswers {
				if answer == strings.ToLower(v) {
					normalAnswer = v
					found = true
					break
				}
			}

			if !found {
				return errPleaseChooseSuggestedAnswer
			}
		} else if item.possibleAnswers != nil && !contains(item.possibleAnswers, answer) {
			return errPleaseChooseSuggestedAnswer
		}
	}

	if item.validateAnswer != nil {
		err := item.validateAnswer(normalAnswer)
		if err != nil {
			return err
		}
	}
	if item.persistAnswer != nil {
		err := item.persistAnswer(normalAnswer, u)
		if err != nil {
			return err
		}
	}
	return nil
}
