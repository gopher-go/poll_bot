package main

import (
	"cloud.google.com/go/datastore"
	"context"
	"time"
)

type datastoreUserDAO struct {
	*datastore.Client
	entityKind string
}

func (d datastoreUserDAO) load(id string) (*storageUser, error) {
	var u storageUser
	if err := d.Get(context.Background(), datastore.NameKey(d.entityKind, id, nil), &u); err != nil {
		return nil, err
	}

	return &u, nil
}

func (d datastoreUserDAO) delete(id string) error {
	return d.Delete(context.Background(), datastore.NameKey(d.entityKind, id, nil))
}

func (d datastoreUserDAO) count() (int, error) {
	return d.Count(context.Background(), datastore.NewQuery(d.entityKind))
}

func (d datastoreUserDAO) save(user *storageUser) error {
	user.CreatedAt = time.Now().UTC()
	_, err := d.Put(context.Background(), datastore.NameKey(d.entityKind, user.ID, nil), user)
	return err
}

func newDatastoreUserDAO(c *datastore.Client, entityKind string) *datastoreUserDAO {
	return &datastoreUserDAO{
		Client:     c,
		entityKind: entityKind,
	}
}
