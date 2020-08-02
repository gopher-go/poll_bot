package main

import (
	"database/sql"
	"github.com/coocood/freecache"
	"math/rand"
	"strings"
	"testing"
	"time"

	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/require"
)

func randomString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String() // E.g. "ExcbsVQs"
	return str
}

func newTestPersistenseStorage() (*persistenseStorage, error) {
	db, err := sql.Open("ramsql", randomString())
	if err != nil {
		return nil, err
	}
	return &persistenseStorage{
		db: db,
	}, nil
}

var (
	testPersistentStorage *persistenseStorage
)

func newTestStorage() (*storage, error) {
	var err error
	testPersistentStorage, err = newTestPersistenseStorage()
	if err != nil {
		return nil, err
	}

	return &storage{
		cache: freecache.NewCache(1024),
		//users:      sync.Map{},
		userDAO: testPersistentStorage,
	}, nil
}

func (s *storage) init() error {
	batch := []string{initDbSQL}

	for _, b := range batch {
		_, err := testPersistentStorage.db.Exec(b)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestMapStorage(t *testing.T) {
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)

	user, err := s.Obtain("12")
	require.NoError(t, err)
	require.Equal(t, user.ID, "12")
	require.Equal(t, user.Properties["age"], "")
	user.Properties["age"] = "16"
	user.Country = "DE"

	err = s.Persist(user)
	require.NoError(t, err)

	user, err = s.Obtain("12")
	require.NoError(t, err)
	require.Equal(t, user.ID, "12")
	require.Equal(t, user.Properties["age"], "16")
	require.Equal(t, user.Country, "DE")

	count, err := s.Count()
	require.NoError(t, err)
	require.Equal(t, count, 1)

	err = s.Persist(user)
	require.NoError(t, err)

	count, err = s.Count()
	require.NoError(t, err)
	require.Equal(t, count, 1)

	// then load persisted
	user, err = s.fromPersisted("12")
	require.NoError(t, err)
	require.Equal(t, user.ID, "12")
	require.Equal(t, user.Properties["age"], "16")
	require.Equal(t, user.Country, "DE")
}
