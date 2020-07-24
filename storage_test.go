package main

import (
	"database/sql"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	_ "github.com/proullon/ramsql/driver"
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

func newTestStorage() (*Storage, error) {
	persistenseStorage, err := newTestPersistenseStorage()
	if err != nil {
		return nil, err
	}

	return &Storage{
		//users:      sync.Map{},
		persistent: persistenseStorage,
	}, nil
}

func (s *Storage) init() error {
	batch := []string{initDbSql}

	for _, b := range batch {
		_, err := s.persistent.db.Exec(b)
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
	require.Equal(t, user.Id, "12")
	require.Equal(t, user.Properties["age"], "")
	user.Properties["age"] = "16"
	user.Country = "DE"
	user.Name = "Georgy"

	err = s.Persist(user)
	require.NoError(t, err)

	user, err = s.Obtain("12")
	require.NoError(t, err)
	require.Equal(t, user.Id, "12")
	require.Equal(t, user.Properties["age"], "16")
	require.Equal(t, user.Country, "DE")
	require.Equal(t, user.Name, "Georgy")

	count, err := s.PersistCount()
	require.NoError(t, err)
	require.Equal(t, count, 1)

	err = s.Persist(user)
	require.NoError(t, err)

	count, err = s.PersistCount()
	require.NoError(t, err)
	require.Equal(t, count, 1)

	// then load persisted
	user, err = s.fromPersisted("12")
	require.NoError(t, err)
	require.Equal(t, user.Id, "12")
	require.Equal(t, user.Properties["age"], "16")
	require.Equal(t, user.Country, "DE")
	require.Equal(t, user.Name, "Georgy")
}

/*
func TestRealStorage(t *testing.T) {
	err := godotenv.Load()
	require.NoError(t, err)

	s, err := newStorage()
	require.NoError(t, err)

	err = s.Clear("12")
	require.NoError(t, err)

	user, err := s.Obtain("12")
	require.NoError(t, err)
	require.Equal(t, user.Id, "12")
	require.Equal(t, user.Properties["age"], "")
	user.Properties["age"] = "16"

	user, err = s.Obtain("12")
	require.NoError(t, err)
	require.Equal(t, user.Id, "12")
	require.Equal(t, user.Properties["age"], "16")

	err = s.Persist(user)
	require.NoError(t, err)

	count, err := s.PersistCount()
	require.NoError(t, err)
	require.Equal(t, count, 1)

	err = s.Persist(user)
	require.NoError(t, err)

	count, err = s.PersistCount()
	require.NoError(t, err)
	require.Equal(t, count, 1)
}
*/
