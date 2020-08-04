package main

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

	ds := datastoreUserDAO{
		dsclient:   c,
		entityKind: "user",
	}

	err = ds.save(&storageUser{ID: "test", Properties: map[string]string{"test": "test"}, MobileCountryCode: 1})
	require.NoError(t, err)

	su, err := ds.load("test")
	require.NoError(t, err)
	su.CreatedAt = time.Time{}
	require.Equal(t, &storageUser{ID: "test", Properties: map[string]string{"test": "test"}, MobileCountryCode: 1}, su)
}

func _TestDatastoreUserDAOUpdate(t *testing.T) {
	os.Setenv("DATASTORE_EMULATOR_HOST", "localhost:8081")
	c, err := datastore.NewClient(context.Background(), "test")
	require.NoError(t, err)

	kind := "user"
	ds := datastoreUserDAO{
		dsclient:   c,
		entityKind: kind,
	}

	err = ds.save(&storageUser{ID: "test", Properties: map[string]string{"test": "test"}, MobileCountryCode: 0})
	require.NoError(t, err)

	fmt.Println(ds.count())
	var users []storageUser
	// .Filter("mcc =", 0) is not working
	g, err := c.GetAll(context.Background(), datastore.NewQuery(ds.entityKind), &users)
	fmt.Println(g)
	fmt.Println(users)
	for _, v := range users {
		err := ds.Update(&v, 2, 3)
		require.NoError(t, err)
	}
	var users2 []storageUser
	g, err = c.GetAll(context.Background(), datastore.NewQuery(ds.entityKind), &users2)

	for _, v := range users2 {
		fmt.Printf("%+v", v)
	}

	require.NoError(t, err)
	su, err := ds.load("test")
	require.NoError(t, err)
	su.CreatedAt = time.Time{}
	require.Equal(t, &storageUser{ID: "test", Properties: map[string]string{"test": "test"}, MobileCountryCode: 2, MobileNetworkCode: 3}, su)
}
