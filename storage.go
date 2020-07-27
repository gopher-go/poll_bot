package main

import (
	"errors"
	"log"
	"time"
)

type userDAO interface {
	load(id string) (*storageUser, error)
	delete(id string) error
	count() (int, error)
	save(user *storageUser) error
}

type answerLog struct {
	UserID      string    `datastore:"user_id"`
	UserContext string    `datastore:"user_context"`
	QuestionID  string    `datastore:"question_id,noindex"`
	Answer      string    `datastore:"answer,noindex"`
	AnswerLevel int       `datastore:"answer_level,noindex"`
	IsValid     bool      `datastore:"is_valid"`
	CreatedAt   time.Time `datastore:"created_at"`
}

type logDAO interface {
	save(al *answerLog) error
}

type storage struct {
	userDAO userDAO
	logDAO  logDAO
}

func newStorage(ud userDAO, ld logDAO) (*storage, error) {
	return &storage{
		logDAO:  ld,
		userDAO: ud,
	}, nil
}

type storageUser struct {
	ID         string
	Country    string
	Level      int
	Properties map[string]string
	Candidate  string
	Context    string
	CreatedAt  time.Time

	isChanged bool
}

func (u *storageUser) validate() error {
	if u.ID == "" {
		return errors.New("Empty user id")
	}
	return nil
}

func (s *storage) LogAnswer(al *answerLog) error {
	if s.logDAO != nil {
		return s.logDAO.save(al)
	}
	return nil
}

func (s *storage) Obtain(id string) (*storageUser, error) {
	if id == "" {
		return nil, errors.New("Unable to obtain empty id")
	}

	persistedUser, err := s.fromPersisted(id)
	if err != nil {
		return nil, err
	}
	if persistedUser != nil {
		persistedUser.isChanged = false
		return persistedUser, persistedUser.validate()
	}

	newUser := &storageUser{
		ID:         id,
		Properties: map[string]string{},
	}

	newUser.isChanged = false

	return newUser, newUser.validate()
}

// internal
func (s *storage) fromPersisted(id string) (*storageUser, error) {
	if s.userDAO == nil {
		return nil, errors.New("persistence not enabled")
	}
	user, err := s.userDAO.load(id)
	if err != nil {
		log.Printf("Unable to load userDAO user %v", err)
		return nil, nil
	}
	return user, nil
}

// Clear removes user from persistent storage
func (s *storage) Clear(id string) error {
	//s.users.Delete(id)

	err := s.userDAO.delete(id)
	if err != nil {
		log.Printf("Unable to clear userDAO user %v", err)
		return nil
	}

	return nil
}

// PersistCount - shows number of users in persistent storage
func (s *storage) PersistCount() (int, error) {
	if s.userDAO == nil {
		return 0, errors.New("persistence not enabled")
	}

	count, err := s.userDAO.count()
	if err != nil {
		log.Printf("Unable to count userDAO user %v", err)
		return 0, nil
	}

	return count, nil
}

// Persist - save the user in storage
func (s *storage) Persist(user *storageUser) error {
	if s.userDAO == nil {
		return errors.New("persistence not enabled")
	}

	err := s.userDAO.save(user)
	if err != nil {
		log.Printf("Unable to save userDAO user %v", err)
		return nil
	}

	return nil
}
