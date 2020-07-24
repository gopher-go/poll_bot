package main

import (
	"errors"
	"log"
)

type Storage struct {
	//users      sync.Map
	persistent *persistenseStorage
}

func newStorage() (*Storage, error) {
	persistent, err := newPersistenseStoragePq()
	if err != nil {
		return nil, err
	}

	return &Storage{
		//users:      sync.Map{},
		persistent: persistent,
	}, nil
}

type StorageUser struct {
	Id      string
	Name    string
	Country string

	Level int

	Properties map[string]string

	Candidate string
	Context   string

	isChanged bool
}

func (u *StorageUser) validate() error {
	if u.Id == "" {
		return errors.New("Empty user id")
	}
	return nil
}

func (s *Storage) Obtain(id string) (*StorageUser, error) {
	if id == "" {
		return nil, errors.New("Unable to obtain empty id")
	}

	//user := s.fromCache(id)
	//if user != nil {
	//	user.isChanged = false
	//	return user, user.validate()
	//}
	persistedUser, err := s.fromPersisted(id)
	if err != nil {
		return nil, err
	}
	if persistedUser != nil {
		persistedUser.isChanged = false
		return persistedUser, persistedUser.validate()
	}

	newUser := &StorageUser{
		Id:         id,
		Properties: map[string]string{},
	}
	//s.users.Store(id, newUser)
	newUser.isChanged = false

	return newUser, newUser.validate()
}

// internal
//func (s *Storage) fromCache(id string) *StorageUser {
//	user, ok := s.users.Load(id)
//	if ok && user != nil {
//		return user.(*StorageUser)
//	}
//	return nil
//}

// internal
func (s *Storage) fromPersisted(id string) (*StorageUser, error) {
	if s.persistent == nil {
		return nil, errors.New("persistence not enabled")
	}
	user, err := s.persistent.load(id)
	if err != nil {
		log.Printf("Unable to load persistent user %v", err)
		return nil, nil
	}
	return user, nil
}

func (s *Storage) Clear(id string) error {
	//s.users.Delete(id)

	err := s.persistent.clear(id)
	if err != nil {
		log.Printf("Unable to clear persistent user %v", err)
		return nil
	}

	return nil
}

func (s *Storage) PersistCount() (int, error) {
	if s.persistent == nil {
		return 0, errors.New("persistence not enabled")
	}

	count, err := s.persistent.count()
	if err != nil {
		log.Printf("Unable to count persistent user %v", err)
		return 0, nil
	}

	return count, nil
}

func (s *Storage) Persist(user *StorageUser) error {
	if s.persistent == nil {
		return errors.New("persistence not enabled")
	}

	//user := s.fromCache(id)
	//if user == nil {
	//	return fmt.Errorf("%v missed in cache", id)
	//}

	err := s.persistent.save(user)
	if err != nil {
		log.Printf("Unable to save persistent user %v", err)
		return nil
	}

	return nil
}
