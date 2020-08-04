package poll_bot

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"time"

	"github.com/coocood/freecache"
)

type userDAO interface {
	load(id string) (*StorageUser, error)
	delete(id string) error
	count() (int, error)
	save(user *StorageUser) error
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

var (
	userCountCacheKey = []byte("user-count")
)

type storage struct {
	cache   *freecache.Cache
	userDAO userDAO
	logDAO  logDAO
}

func newStorage(ud userDAO, ld logDAO) (*storage, error) {
	return &storage{
		cache:   freecache.NewCache(512 * 1024 * 1024),
		logDAO:  ld,
		userDAO: ud,
	}, nil
}

type StorageUser struct {
	ID                string
	Country           string
	Language          string
	MobileCountryCode int
	MobileNetworkCode int
	Level             int
	Properties        map[string]string
	Candidate         string
	Context           string
	CreatedAt         time.Time

	isChanged bool
}

func (u *StorageUser) validate() error {
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

func (s *storage) Obtain(id string) (*StorageUser, error) {
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

	newUser := &StorageUser{
		ID:         id,
		Properties: map[string]string{},
	}

	newUser.isChanged = false

	return newUser, newUser.validate()
}

// internal
func (s *storage) fromPersisted(id string) (*StorageUser, error) {
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

// Count - shows number of users in persistent storage
func (s *storage) Count() (int, error) {
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

// CountCached - get the number of users using cache
func (s *storage) CountCached() (count int, err error) {
	countBytes, err := s.cache.Get(userCountCacheKey)
	if err == nil {
		err = gob.NewDecoder(bytes.NewBuffer(countBytes)).Decode(&count)
		if err == nil {
			return
		}
	}

	count, err = s.Count()
	if err == nil {
		var countBytesBuff bytes.Buffer
		err = gob.NewEncoder(&countBytesBuff).Encode(count)
		if err == nil {
			_ = s.cache.Set(userCountCacheKey, countBytesBuff.Bytes(), 5)
		}
	}

	return
}

// Persist - save the user in storage
func (s *storage) Persist(user *StorageUser) error {
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
