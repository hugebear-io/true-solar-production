package main

import (
	"context"
	"fmt"

	"github.com/hugebear-io/true-solar-production/config"
	"github.com/hugebear-io/true-solar-production/constant"
	"github.com/hugebear-io/true-solar-production/infra"
	"github.com/hugebear-io/true-solar-production/util"
	"github.com/olivere/elastic/v7"
)

func init() {
	config.InitConfig()
}

func init() {
	util.SetTimezone()
}

func main() {
	e, err := infra.NewElasticsearch()
	if err != nil {
		panic(err)
	}

	compositeAggregation := elastic.NewCompositeAggregation().
		Size(10000).
		Sources(elastic.NewCompositeAggregationDateHistogramValuesSource("date").Field("@timestamp").CalendarInterval("day").Format("yyyy-MM-dd"),
			elastic.NewCompositeAggregationTermsValuesSource("vendor_type").Field("vendor_type.keyword"),
			elastic.NewCompositeAggregationTermsValuesSource("id").Field("id.keyword")).
		SubAggregation("max_daily", elastic.NewMaxAggregation().Field("daily_production")).
		SubAggregation("avg_capacity", elastic.NewAvgAggregation().Field("installed_capacity")).
		SubAggregation("threshold_percentage", elastic.NewBucketScriptAggregation().
			BucketsPathsMap(map[string]string{"capacity": "avg_capacity"}).
			Script(elastic.NewScript("params.capacity * params.efficiency_factor * params.focus_hour * params.threshold_percentage").
				Params(map[string]interface{}{
					"efficiency_factor":    0.8,
					"focus_hour":           5,
					"threshold_percentage": 0.6,
				}))).
		SubAggregation("under_threshold", elastic.NewBucketSelectorAggregation().
			BucketsPathsMap(map[string]string{"threshold": "threshold_percentage", "daily": "max_daily"}).
			Script(elastic.NewScript("params.daily <= params.threshold"))).
		SubAggregation("hits", elastic.NewTopHitsAggregation().
			Size(1).
			FetchSourceContext(
				elastic.NewFetchSourceContext(true).Include(
					"id", "name", "vendor_type", "node_type", "ac_phase", "plant_status",
					"area", "site_id", "site_city_code", "site_city_name", "installed_capacity",
				)))

	searchQuery := e.Search("solarcell*").
		Size(0).
		Query(elastic.NewBoolQuery().Must(
			elastic.NewMatchQuery("data_type", constant.DATA_TYPE_PLANT),
			elastic.NewRangeQuery("@timestamp").Gte("now-7d/d").Lte("now-1d/d"),
		)).
		Aggregation("performance_alarm", compositeAggregation)

	items := make([]*elastic.AggregationBucketCompositeItem, 0)
	result, err := searchQuery.Pretty(true).Do(context.Background())
	if err != nil {
		panic(err)
	}

	if result.Aggregations == nil {
		panic("no aggregation")
	}

	performanceAlarm, found := result.Aggregations.Composite("performance_alarm")
	if !found {
		panic("no performance_alarm aggregation")
	}

	items = append(items, performanceAlarm.Buckets...)
	util.PrintJSON(map[string]interface{}{"result": items, "after_key": performanceAlarm.AfterKey})

	if len(performanceAlarm.AfterKey) > 0 {
		afterKey := performanceAlarm.AfterKey

		for {
			query := e.Search("solarcell*").
				Size(0).
				Query(elastic.NewBoolQuery().Must(
					elastic.NewMatchQuery("data_type", constant.DATA_TYPE_PLANT),
					elastic.NewRangeQuery("@timestamp").Gte("now-7d/d").Lte("now-1d/d"),
				)).
				Aggregation("performance_alarm", compositeAggregation.AggregateAfter(afterKey))

			result, err := query.Pretty(true).Do(context.Background())
			if err != nil {
				panic(err)
			}

			if result.Aggregations == nil {
				panic("no aggregation")
			}

			performanceAlarm, found := result.Aggregations.Composite("performance_alarm")
			if !found {
				panic("no performance_alarm aggregation")
			}

			items = append(items, performanceAlarm.Buckets...)

			if len(performanceAlarm.AfterKey) == 0 {
				break
			}

			afterKey = performanceAlarm.AfterKey
		}
	}

	fmt.Println("=====================================")
	fmt.Println("=====================================")
	fmt.Println("=====================================")

	util.PrintJSON(map[string]interface{}{"result": items, "after_key": performanceAlarm.AfterKey})
}
