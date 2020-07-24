package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
)

type redisUserDAO struct {
	client *redis.Client
}

func newRedisUserDAO(client *redis.Client) *redisUserDAO {
	return &redisUserDAO{
		client: client,
	}
}

func (r redisUserDAO) load(id string) (*StorageUser, error) {
	repl, err := r.client.Get(context.Background(), id).Bytes()
	if err != nil {
		return nil, err
	}

	var u StorageUser
	if err := json.Unmarshal(repl, &u); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r redisUserDAO) delete(id string) error {
	return r.client.Del(context.Background(), id).Err()
}

func (r redisUserDAO) count() (int, error) {
	count, err := r.client.DBSize(context.Background()).Result()
	return int(count), err
}

func (r redisUserDAO) save(user *StorageUser) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.client.Set(context.Background(), user.Id, userBytes, 0).Err()
}
