package main

import (
	"context"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
)

type dsclient interface {
	Count(ctx context.Context, q *datastore.Query) (n int, err error)
	Delete(ctx context.Context, key *datastore.Key) error
	Get(ctx context.Context, key *datastore.Key, dst interface{}) (err error)
	Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error)
}

type datastoreUserDAO struct {
	dsclient
	entityKind string
}

func (u *storageUser) Load(properties []datastore.Property) error {
	propMap := map[string]datastore.Property{}
	for i := range properties {
		propMap[properties[i].Name] = properties[i]
	}

	u.ID = propMap["id"].Value.(string)
	u.Country = propMap["country"].Value.(string)
	u.Candidate = propMap["candidate"].Value.(string)
	u.Context = propMap["context"].Value.(string)
	u.Level = int(propMap["level"].Value.(int64))
	u.CreatedAt = propMap["created_at"].Value.(time.Time)

	u.Properties = map[string]string{}
	for i := range properties {
		if strings.HasPrefix(properties[i].Name, "property.") {
			u.Properties[properties[i].Name[9:]] = properties[i].Value.(string)
		}
	}

	return nil
}

func (u *storageUser) Save() ([]datastore.Property, error) {
	var props []datastore.Property

	props = append(props, datastore.Property{Name: "id", Value: u.ID})
	props = append(props, datastore.Property{Name: "created_at", Value: u.CreatedAt})
	props = append(props, datastore.Property{Name: "context", Value: u.Context})
	props = append(props, datastore.Property{Name: "level", Value: u.Level, NoIndex: true})
	props = append(props, datastore.Property{Name: "candidate", Value: u.Candidate})
	props = append(props, datastore.Property{Name: "country", Value: u.Country})

	for pn, p := range u.Properties {
		props = append(props, datastore.Property{Name: "property." + pn, Value: p, NoIndex: true})
	}

	return props, nil
}

func (d datastoreUserDAO) load(id string) (*storageUser, error) {
	var u storageUser
	if err := d.Get(context.Background(), datastore.NameKey(d.entityKind, id, nil), &u); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil
		}
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
		dsclient:   c,
		entityKind: entityKind,
	}
}
