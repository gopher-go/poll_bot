package main

import (
	"context"
	"strconv"

	"cloud.google.com/go/datastore"
)

type datastoreLogDAO struct {
	*datastore.Client
	entityKind string
}

func newDatastoreLogDAO(c *datastore.Client, entityKind string) *datastoreLogDAO {
	return &datastoreLogDAO{
		Client:     c,
		entityKind: entityKind,
	}
}

func datastoreAnswerLogRowKey(al answerLog) string {
	return al.UserID + "#" + al.UserContext + "#" + strconv.FormatInt(al.CreatedAt.UnixNano(), 16)
}

func (d datastoreLogDAO) save(al answerLog) error {
	_, err := d.Put(context.Background(), datastore.NameKey(d.entityKind, datastoreAnswerLogRowKey(al), nil), al)
	return err
}
