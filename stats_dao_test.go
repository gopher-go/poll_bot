package poll_bot

import (
	"context"
	"testing"

	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/require"
)

func _TestStatsDaoCountByField(t *testing.T) {
	c, err := elastic.NewSimpleClient(elastic.SetURL("https://localhoost:9200"))
	require.NoError(t, err)
	td := newStatsDao(c)

	agg, err := td.CountByFieldCached(context.Background(), residenceLocationType, residenceLocation)
	require.NoError(t, err)
	agg, err = td.CountByFieldCached(context.Background(), residenceLocationType, residenceLocation)
	require.NoError(t, err)
	_ = agg
}
