package poll_bot

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/coocood/freecache"
	"github.com/olivere/elastic/v7"
)

const (
	residenceLocation     = "property.residence_location.keyword"
	residenceLocationType = "property.residence_location_type.keyword"
)

type statsDao struct {
	cache  *freecache.Cache
	client *elastic.Client
}

func newStatsDao(client *elastic.Client) *statsDao {
	return &statsDao{
		cache:  freecache.NewCache(512 * 1024 * 1024),
		client: client,
	}
}

func newTermAggregations(fields []string) (aggs map[string]*elastic.TermsAggregation) {
	aggs = map[string]*elastic.TermsAggregation{}
	for _, f := range fields {
		aggs[string(f)] = elastic.NewTermsAggregation().
			Field(string(f)).
			Order("_count", false).
			Size(12)

	}
	return
}

// CountByField - counts by field
func (sd *statsDao) CountByField(ctx context.Context, field ...string) (map[string]map[string]int, error) {
	query := sd.client.Search("users").
		Query(elastic.NewBoolQuery().Must(
			elastic.NewBoolQuery().MustNot(elastic.NewTermQuery("property.residence_location.keyword", "Проживаю за пределами РБ")),
			elastic.NewBoolQuery().MustNot(elastic.NewTermQuery("property.residence_location_type.keyword", "Проживаю за пределами РБ")),
			elastic.NewBoolQuery().MustNot(elastic.NewTermQuery("property.age.keyword", "младше 18")),
			elastic.NewBoolQuery().MustNot(elastic.NewTermQuery("candidate.keyword", "")),
			elastic.NewTermQuery("country.keyword", "BY"),
			elastic.NewTermQuery("property.isBelarus.keyword", "Да"),
		)).
		Size(0)

	aggs := newTermAggregations(field)
	for f, a := range aggs {
		query.Aggregation(f, a)
	}

	searchResult, err := query.Do(ctx)
	if err != nil {
		return nil, err
	}

	result := map[string]map[string]int{}
	for _, f := range field {
		result[string(f)] = map[string]int{}

		t, ok := searchResult.Aggregations.Terms(string(f))
		if !ok {
			return result, nil
		}

		for i := range t.Buckets {
			if t.Buckets[i] != nil {
				k, _ := t.Buckets[i].Key.(string)
				result[string(f)][k] = int(t.Buckets[i].DocCount)
			}
		}
	}

	return result, nil
}

//  CountByFieldCached counts by field returns cached results if available
func (sd *statsDao) CountByFieldCached(ctx context.Context, field ...string) (result map[string]map[string]int, err error) {
	cacheKey := []byte(strings.Join(field, "$"))
	cachedResultBytes, err := sd.cache.Get(cacheKey)
	if err == nil {
		if err := json.Unmarshal(cachedResultBytes, &result); err == nil {
			return result, nil
		}
	}

	result, err = sd.CountByField(ctx, field...)
	if err == nil {
		resultBytes, err := json.Marshal(result)
		if err == nil {
			_ = sd.cache.Set(cacheKey, resultBytes, 30)
		}
	}

	return
}
