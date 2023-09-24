package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hugebear-io/true-solar-production/config"
	"github.com/hugebear-io/true-solar-production/constant"
	"github.com/olivere/elastic/v7"
)

type SolarRepo interface {
	BulkIndex(index string, docs []interface{}) error
	GetPlantDailyProduction(start, end *time.Time) ([]*elastic.AggregationBucketCompositeItem, error)
}

type solarRepo struct {
	elastic *elastic.Client
}

func NewSolarRepo(elastic *elastic.Client) SolarRepo {
	return &solarRepo{elastic: elastic}
}

func (r *solarRepo) searchIndex() *elastic.SearchService {
	conf := config.GetConfig().Elastic
	index := fmt.Sprintf("%v*", conf.SolarIndex)
	return r.elastic.Search(index)
}

func (r *solarRepo) createIndexIfNotExist(index string) error {
	ctx := context.Background()
	if exist, err := r.elastic.IndexExists(index).Do(ctx); err != nil {
		if !exist {
			result, err := r.elastic.CreateIndex(index).Do(ctx)
			if err != nil {
				return err
			}

			if !result.Acknowledged {
				return errors.New("elasticsearch did not acknowledge")
			}
		}
	}

	return nil
}

// |=> Implementation
func (r *solarRepo) BulkIndex(index string, docs []interface{}) error {
	if err := r.createIndexIfNotExist(index); err != nil {
		return err
	}

	bulk := r.elastic.Bulk()
	for _, doc := range docs {
		bulk.Add(elastic.NewBulkIndexRequest().Index(index).Doc(doc))
	}

	ctx := context.Background()
	if _, err := bulk.Do(ctx); err != nil {
		return err
	}

	return nil
}

func (r *solarRepo) GetPlantDailyProduction(start, end *time.Time) ([]*elastic.AggregationBucketCompositeItem, error) {
	ctx := context.Background()
	items := make([]*elastic.AggregationBucketCompositeItem, 0)

	// create [composite aggregation]
	compositeAggregation := elastic.NewCompositeAggregation().
		Size(10000).
		Sources(
			elastic.NewCompositeAggregationDateHistogramValuesSource("date").Field("@timestamp").CalendarInterval("day").Format("yyyy-MM-dd"),
			elastic.NewCompositeAggregationTermsValuesSource("vendor_type").Field("vendor_type.keyword"),
			elastic.NewCompositeAggregationTermsValuesSource("area").Field("area.keyword"),
			elastic.NewCompositeAggregationTermsValuesSource("name").Field("name.keyword"),
		)

	// assign [max_aggregation] into composite aggregation
	compositeAggregation = compositeAggregation.
		SubAggregation("installed_capacity", elastic.NewMaxAggregation().Field("installed_capacity")).
		SubAggregation("monthly_production", elastic.NewMaxAggregation().Field("monthly_production")).
		SubAggregation("daily_production", elastic.NewMaxAggregation().Field("daily_production"))

	// assign [hit_aggregation] into composite aggregation
	compositeAggregation = compositeAggregation.
		SubAggregation(
			"hits",
			elastic.NewTopHitsAggregation().
				Size(1).
				FetchSourceContext(elastic.NewFetchSourceContext(true).Include("lat", "lng")),
		)

	// assign [bucket_script_aggregation] into composite aggregation
	// |=> target
	const targetScript = "if (params.installed_capacity == 0 || params.daily_production == 0 ) { return 0 } else { (params.installed_capacity*5*0.8) }"

	compositeAggregation = compositeAggregation.SubAggregation(
		"target",
		elastic.NewBucketScriptAggregation().
			BucketsPathsMap(map[string]string{"installed_capacity": "installed_capacity", "daily_production": "daily_production"}).
			Script(elastic.NewScript(targetScript)),
	)

	// |=> production_to_target
	const productionToTargetScript = "if (params.installed_capacity == 0 || params.daily_production == 0 ) { return 0 } else { (params.daily_production/(params.installed_capacity*5*0.8))*100 }"

	compositeAggregation = compositeAggregation.SubAggregation(
		"production_to_target",
		elastic.NewBucketScriptAggregation().
			BucketsPathsMap(map[string]string{"installed_capacity": "installed_capacity", "daily_production": "daily_production"}).
			Script(elastic.NewScript(productionToTargetScript)),
	)

	query := r.searchIndex().
		Size(0).
		Query(
			elastic.NewBoolQuery().Must(
				elastic.NewMatchQuery("data_type", constant.DATA_TYPE_PLANT),
				elastic.NewRangeQuery("@timestamp").Gte(start).Lt(end),
			),
		).Aggregation("production", compositeAggregation)

	result, err := query.Pretty(true).Do(ctx)
	if err != nil {
		return nil, err
	}

	if result.Aggregations == nil {
		return nil, errors.New("cannot get result aggregation")
	}

	productions, found := result.Aggregations.Composite("production")
	if !found {
		return nil, errors.New("cannot get result composite performance alarm")
	}

	items = append(items, productions.Buckets...)
	return items, nil
}
