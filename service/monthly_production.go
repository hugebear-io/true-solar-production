package service

import (
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/hugebear-io/true-solar-production/config"
	"github.com/hugebear-io/true-solar-production/constant"
	"github.com/hugebear-io/true-solar-production/logger"
	"github.com/hugebear-io/true-solar-production/model"
	"github.com/hugebear-io/true-solar-production/repo"
)

type MonthlyProductionService interface {
	Run(start, end *time.Time) error
}

type monthlyProductionService struct {
	solarRepo repo.SolarRepo
	logger    logger.Logger
}

func NewMonthlyProductionService(solarRepo repo.SolarRepo, logger logger.Logger) MonthlyProductionService {
	return &monthlyProductionService{solarRepo: solarRepo, logger: logger}
}

func (s monthlyProductionService) Run(start, end *time.Time) error {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Warnf("[%v]MonthlyProduction.Run(): %v", start.Format(constant.YEAR_MONTH), r)
		}
	}()

	documents, err := s.generateDocuments(start, end)
	if err != nil {
		s.logger.Errorf("[%v]MonthlyProduction.Run(): %v", start.Format(constant.YEAR_MONTH), err)
		return err
	}

	if len(documents) == 0 {
		s.logger.Errorf("[%v]MonthlyProduction.Run(): %v", start.Format(constant.YEAR_MONTH), "documents is empty")
		return nil
	}

	conf := config.GetConfig().Elastic
	index := fmt.Sprintf("%v_%v", conf.MonthlyProductionIndex, start.Format("200601"))
	if err := s.solarRepo.BulkIndex(index, documents); err != nil {
		s.logger.Errorf("[%v]MonthlyProduction.Run(): %v", start.Format(constant.YEAR_MONTH), err)
		return err
	}

	s.logger.Infof("MonthlyProduction.Run(): bulked new index %q", index)
	return nil
}

func (s monthlyProductionService) generateDocuments(start, end *time.Time) ([]interface{}, error) {
	data, err := s.solarRepo.GetPlantMonthlyProduction(start, end)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	var count int
	var size = len(data)
	documents := make([]interface{}, 0)
	for _, item := range data {
		if item == nil {
			continue
		}

		if len(item.Key) == 0 {
			continue
		}

		doc := model.MonthlyProductionDocument{}

		// took data from key
		if val, ok := item.Key["date"].(string); ok {
			date, _ := time.Parse(constant.YEAR_MONTH_DAY, val)
			doc.SetDate(&date)
		}

		if val, ok := item.Key["vendor_type"].(string); ok {
			doc.SetVendorType(val)
		}

		if val, ok := item.Key["area"].(string); ok {
			doc.SetArea(val)
		}

		if val, ok := item.Key["name"].(string); ok {
			doc.SetSiteName(val)
		}

		// took data from max_aggregation
		if val, ok := item.Aggregations.Max("installed_capacity"); ok {
			doc.SetInstalledCapacity(val.Value)
		}

		if val, ok := item.Aggregations.Max("monthly_production"); ok {
			doc.SetMonthlyProduction(val.Value)
		}

		// took data from bucket script
		if val, ok := item.BucketScript("target"); ok {
			doc.SetTarget(val.Value)
		}

		if val, ok := item.BucketScript("production_to_target"); ok {
			doc.SetProductionToTarget(val.Value)
		}
		doc.SetCriteria(doc.ProductionToTarget)

		// took data from hits
		if hits, ok := item.TopHits("hits"); ok {
			if hits.Hits != nil {
				if len(hits.Hits.Hits) == 1 {
					hit := hits.Hits.Hits[0]
					source := make(map[string]float64)
					if err := sonic.Unmarshal(hit.Source, &source); err == nil {
						if val, ok := source["lat"]; ok {
							doc.SetLatitude(&val)
						}

						if val, ok := source["lng"]; ok {
							doc.SetLongitude(&val)
						}
					}
				}
			}
		}

		count += 1
		s.logger.Infof("[%v/%v] generateDocument of %v", count, size, start.Format(constant.YEAR_MONTH))
		documents = append(documents, doc)
	}

	return documents, nil
}
