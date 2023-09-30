package handler

import (
	"time"

	"github.com/gammazero/workerpool"
	"github.com/hugebear-io/true-solar-production/constant"
	"github.com/hugebear-io/true-solar-production/infra"
	"github.com/hugebear-io/true-solar-production/logger"
	"github.com/hugebear-io/true-solar-production/repo"
	"github.com/hugebear-io/true-solar-production/service"
)

type MonthlyProductionHandler struct{}

func NewMonthlyProductionHandler() *MonthlyProductionHandler {
	return &MonthlyProductionHandler{}
}

func (h MonthlyProductionHandler) Run() {
	pool := workerpool.New(5)
	currentMonth := time.January
	endDate := time.Now()

	for {
		start := time.Date(2023, currentMonth, 1, 0, 0, 0, 0, time.Local)
		end := time.Date(2023, currentMonth+1, 1, 0, 0, 0, 0, time.Local)

		pool.Submit(h.run(&start, &end))
		if endDate.Month() == currentMonth {
			break
		} else {
			currentMonth += 1
		}
	}

	pool.StopWait()
}

func (h MonthlyProductionHandler) run(start, end *time.Time) func() {
	return func() {
		logger := logger.NewLogger(
			&logger.LoggerOption{
				LogName:     "logs/monthly_production.log",
				LogSize:     1024,
				LogAge:      90,
				LogBackup:   1,
				LogCompress: false,
				LogLevel:    logger.LOG_LEVEL_DEBUG,
				SkipCaller:  1,
			},
		)
		defer logger.Close()

		elastic, err := infra.NewElasticsearch()
		if err != nil {
			logger.Errorf("[%v]Failed to connect to elasticsearch", start.Format(constant.YEAR_MONTH))
			return
		}

		masterSiteRepo, err := repo.NewMasterSiteRepo()
		if err != nil {
			logger.Error(err)
			return
		}

		solarRepo := repo.NewSolarRepo(elastic)
		serv := service.NewMonthlyProductionService(solarRepo, masterSiteRepo, logger)
		if err := serv.Run(start, end); err != nil {
			logger.Error(err)
		}
	}
}
