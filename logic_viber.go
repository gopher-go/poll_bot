package main

import (
	"errors"
	"fmt"
	"strings"
)

func knownEvent(c *ViberCallback) bool {
	return c.Event == "message" ||
		c.Event == "delivered" ||
		c.Event == "seen" ||
		c.Event == "subscribed" ||
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

func generateReplyFor(p poll, s *Storage, c *ViberCallback) (*viberReply, error) {
	if !knownEvent(c) {
		return nil, fmt.Errorf("Unknown message %v", c.Event)
	}

	if strings.ToLower(c.Message.Text) == "clear" {
		err := s.Clear(c.User.Id)
		return &viberReply{text: fmt.Sprintf("Your storage cleared with %v", err)}, nil
	}

	storageUser, err := s.Obtain(c.User.Id)
	if err != nil {
		return nil, err
	}

	defer func() {
		if storageUser.isChanged {
			_ = s.Persist(storageUser.Id)
		}
	}()

	if storageUser.Name == "" && c.User.Name != "" {
		storageUser.Name = c.User.Name
		storageUser.isChanged = true
	}

	if storageUser.Country == "" && c.User.Country != "" {
		storageUser.Country = c.User.Country
		storageUser.isChanged = true
	}

	if c.Event == "message" {
		err := analyseAnswer(p, storageUser, c)
		if err != nil {
			reply := getViberReplyForLevel(p, storageUser, storageUser.Level, c)
			reply.text = err.Error() + " " + reply.text
			return reply, nil
		}
		storageUser.Level++
		storageUser.isChanged = true
		if storageUser.Level >= p.size {
			storageUser.isChanged = true

			totalCount, err := s.PersistCount()
			if err != nil {
				return nil, err
			}
			text := fmt.Sprintf("Спасибо за голосование! Уже проголосовало %d человек", totalCount)
			return &viberReply{text: text}, nil
		}

		reply := getViberReplyForLevel(p, storageUser, storageUser.Level, c)
		return reply, nil
	}

	if !storageUser.ConversationStarted {
		reply := getViberReplyForLevel(p, storageUser, storageUser.Level, c)
		storageUser.ConversationStarted = true
		storageUser.isChanged = true
		return reply, nil
	}

	return nil, nil
}

func getViberReplyForLevel(p poll, u *StorageUser, level int, c *ViberCallback) *viberReply {
	item := p.getLevel(level)
	reply := viberReply{text: fmt.Sprintf("Непонятно. Нет уровня %v в вопросах", level)}
	if item != nil {
		reply.text = item.question(u, c)
		reply.options = item.possibleAnswers
	}
	return &reply
}

func analyseAnswer(p poll, u *StorageUser, c *ViberCallback) error {
	item := p.getLevel(u.Level)
	if item == nil {
		return nil
	}

	answer := c.Message.Text
	normalAnswer := answer
	if !caseSensitive {
		answer = strings.ToLower(answer)
		found := false
		for _, v := range item.possibleAnswers {
			if answer == strings.ToLower(v) {
				normalAnswer = v
				found = true
				break
			}
		}
		if !found {
			return errors.New("Пожалуйста выберите предложенный ответ.")
		}
	} else if item.possibleAnswers != nil && !contains(item.possibleAnswers, answer) {
		return errors.New("Пожалуйста выберите предложенный ответ.")
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
