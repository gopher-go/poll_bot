package poll_bot

import (
	"context"
	"strconv"

	"cloud.google.com/go/datastore"
)

type datastoreLogDAO struct {
	*datastore.Client
	EntityKind string
}

func newDatastoreLogDAO(c *datastore.Client, EntityKind string) *datastoreLogDAO {
	return &datastoreLogDAO{
		Client:     c,
		EntityKind: EntityKind,
	}
}

func datastoreAnswerLogRowKey(al *answerLog) string {
	return al.UserID + "#" + al.UserContext + "#" + strconv.FormatInt(al.CreatedAt.UnixNano(), 16)
}

func (d datastoreLogDAO) save(al *answerLog) error {
	_, err := d.Put(context.Background(), datastore.NameKey(d.EntityKind, datastoreAnswerLogRowKey(al), nil), al)
	return err
}
