package poll_bot

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/stretchr/testify/require"
)

func _TestDatastoreUserDAOSave(t *testing.T) {
	os.Setenv("DATASTORE_EMULATOR_HOST", "localhost:8081")
	c, err := datastore.NewClient(context.Background(), "test")
	require.NoError(t, err)

	ds := DatastoreUserDAO{
		DSClient:   c,
		EntityKind: "user",
	}

	err = ds.save(&StorageUser{ID: "test", Properties: map[string]string{"test": "test"}, MobileCountryCode: 1})
	require.NoError(t, err)

	su, err := ds.load("test")
	require.NoError(t, err)
	su.CreatedAt = time.Time{}
	require.Equal(t, &StorageUser{ID: "test", Properties: map[string]string{"test": "test"}, MobileCountryCode: 1}, su)
}

func TestDatastoreUserDAOUpdate(t *testing.T) {
	os.Setenv("DATASTORE_EMULATOR_HOST", "localhost:8081")
	c, err := datastore.NewClient(context.Background(), "test")
	require.NoError(t, err)

	kind := "user"
	ds := DatastoreUserDAO{
		DSClient:   c,
		EntityKind: kind,
	}

	err = ds.save(&StorageUser{ID: "test", Country: "Bel", Context: "IP", Language: "e", Candidate: "Тих", Properties: map[string]string{"test": "test"}, MobileCountryCode: 0})

	err = ds.save(&StorageUser{ID: "test2", Country: "Bel2", Language: "e", Candidate: "Тих", Properties: map[string]string{"test": "test"}, MobileCountryCode: 0})
	require.NoError(t, err)

	fmt.Println(ds.count())
	var users []StorageUser
	// .Filter("mcc =", 0) is not working
	g, err := c.GetAll(context.Background(), datastore.NewQuery(ds.EntityKind), &users)
	fmt.Println(g)
	fmt.Println(users)
	for _, v := range users {
		err := ds.Update(&v, 2, 3)
		require.NoError(t, err)
	}
	var users2 []StorageUser
	g, err = c.GetAll(context.Background(), datastore.NewQuery(ds.EntityKind), &users2)

	for _, v := range users2 {
		fmt.Printf("%+v", v)
	}

	require.NoError(t, err)
	su, err := ds.load("test")
	require.NoError(t, err)
	su.CreatedAt = time.Time{}
	require.Equal(t, &StorageUser{ID: "test", Properties: map[string]string{"test": "test"}, Context: "IP", Candidate: "Тих", Country: "Bel", Language: "e", MobileCountryCode: 2, MobileNetworkCode: 3}, su)

	su, err = ds.load("test2")
	require.NoError(t, err)
	su.CreatedAt = time.Time{}
	require.Equal(t, &StorageUser{ID: "test2", Properties: map[string]string{"test": "test"}, Candidate: "Тих", Country: "Bel2", Language: "e", MobileCountryCode: 2, MobileNetworkCode: 3}, su)
}
