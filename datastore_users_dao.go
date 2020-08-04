package poll_bot

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
)

type DSClient interface {
	Count(ctx context.Context, q *datastore.Query) (n int, err error)
	Delete(ctx context.Context, key *datastore.Key) error
	Get(ctx context.Context, key *datastore.Key, dst interface{}) (err error)
	Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error)
	NewTransaction(ctx context.Context, opts ...datastore.TransactionOption) (t *datastore.Transaction, err error)
	GetAll(ctx context.Context, q *datastore.Query, dst interface{}) (keys []*datastore.Key, err error)
}

type DatastoreUserDAO struct {
	DSClient
	EntityKind string
}

func getInt(i interface{}) int {
	iv := reflect.ValueOf(i)

	if !iv.IsValid() {
		return 0
	}

	return int(iv.Int())
}

func getString(i interface{}) string {
	iv := reflect.ValueOf(i)
	if !iv.IsValid() {
		return ""
	}

	return iv.String()
}

func (u *StorageUser) Load(properties []datastore.Property) error {
	propMap := map[string]datastore.Property{}
	for i := range properties {
		propMap[properties[i].Name] = properties[i]
	}

	u.ID = getString(propMap["id"].Value)
	u.Country = getString(propMap["country"].Value)
	u.Candidate = getString(propMap["candidate"].Value)
	u.Context = getString(propMap["context"].Value)
	u.Level = getInt(propMap["level"].Value)
	u.MobileCountryCode = getInt(propMap["mcc"].Value)
	u.MobileNetworkCode = getInt(propMap["mnc"].Value)
	u.Language = getString(propMap["language"].Value)
	u.CreatedAt = propMap["created_at"].Value.(time.Time)

	u.Properties = map[string]string{}
	for i := range properties {
		if strings.HasPrefix(properties[i].Name, "property.") {
			u.Properties[properties[i].Name[9:]] = getString(properties[i].Value)
		}
	}

	return nil
}

func (u *StorageUser) Save() ([]datastore.Property, error) {
	var props []datastore.Property

	props = append(props, datastore.Property{Name: "id", Value: u.ID})
	props = append(props, datastore.Property{Name: "created_at", Value: u.CreatedAt})
	props = append(props, datastore.Property{Name: "context", Value: u.Context})
	props = append(props, datastore.Property{Name: "level", Value: u.Level, NoIndex: true})
	props = append(props, datastore.Property{Name: "candidate", Value: u.Candidate})
	props = append(props, datastore.Property{Name: "country", Value: u.Country})
	props = append(props, datastore.Property{Name: "language", Value: u.Language, NoIndex: true})
	props = append(props, datastore.Property{Name: "mcc", Value: u.MobileCountryCode, NoIndex: true})
	props = append(props, datastore.Property{Name: "mnc", Value: u.MobileNetworkCode, NoIndex: true})

	for pn, p := range u.Properties {
		props = append(props, datastore.Property{Name: "property." + pn, Value: p, NoIndex: true})
	}

	return props, nil
}

func (d DatastoreUserDAO) load(id string) (*StorageUser, error) {
	var u StorageUser
	if err := d.Get(context.Background(), datastore.NameKey(d.EntityKind, id, nil), &u); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (d DatastoreUserDAO) delete(id string) error {
	return d.Delete(context.Background(), datastore.NameKey(d.EntityKind, id, nil))
}

func (d DatastoreUserDAO) count() (int, error) {
	return d.Count(context.Background(), datastore.NewQuery(d.EntityKind))
}

func (d DatastoreUserDAO) save(user *StorageUser) error {
	user.CreatedAt = time.Now().UTC()
	_, err := d.Put(context.Background(), datastore.NameKey(d.EntityKind, user.ID, nil), user)
	return err
}

func (d DatastoreUserDAO) Update(user *StorageUser, MCC, MNC int) error {
	tx, err := d.DSClient.NewTransaction(context.Background())
	if err != nil {
		return fmt.Errorf("client.NewTransaction: %v", err)
	}
	var userGet StorageUser
	userKey := datastore.NameKey(d.EntityKind, user.ID, nil)
	if err := tx.Get(userKey, &userGet); err != nil {
		return fmt.Errorf("tx.Get: %v", err)
	}
	userGet.MobileCountryCode = MCC
	userGet.MobileNetworkCode = MNC
	if _, err := tx.Put(userKey, &userGet); err != nil {
		return fmt.Errorf("tx.Put: %v", err)
	}
	if _, err := tx.Commit(); err != nil {
		return fmt.Errorf("tx.Commit: %v", err)
	}
	return err
}

func NewDatastoreUserDAO(c *datastore.Client, EntityKind string) *DatastoreUserDAO {
	return &DatastoreUserDAO{
		DSClient:   c,
		EntityKind: EntityKind,
	}
}
