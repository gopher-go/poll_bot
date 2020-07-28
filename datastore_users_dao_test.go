package main

import (
	"context"
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
