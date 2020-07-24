package main

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
)

func TestRedisUserDAO(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	defer mr.Close()

	rud := newRedisUserDAO(redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	}))

	// save
	err = rud.save(&StorageUser{
		Id:         "1234",
		Name:       "test-name",
		Country:    "-",
		Level:      2,
		Properties: map[string]string{"key": "value"},
		Candidate:  "test-cadidate",
		Context:    "sms",
		isChanged:  true,
	})
	require.NoError(t, err)

	// load existing
	u, err := rud.load("1234")
	require.NoError(t, err)
	require.Equal(t, &StorageUser{
		Id:         "1234",
		Name:       "test-name",
		Country:    "-",
		Level:      2,
		Properties: map[string]string{"key": "value"},
		Candidate:  "test-cadidate",
		Context:    "sms",
		isChanged:  false,
	}, u)

	// count
	uc, err := rud.count()
	require.NoError(t, err)
	require.Equal(t, 1, uc)

	// delete
	err = rud.delete("1234")
	require.NoError(t, err)
}
