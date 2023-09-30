package handler

import (
	"math"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/hugebear-io/true-solar-production/constant"
	"github.com/hugebear-io/true-solar-production/infra"
	"github.com/hugebear-io/true-solar-production/logger"
	"github.com/hugebear-io/true-solar-production/repo"
	"github.com/hugebear-io/true-solar-production/service"
)

type DailyProductionHandler struct{}

func NewDailyProductionHandler() *DailyProductionHandler {
	return &DailyProductionHandler{}
}

func (h DailyProductionHandler) Run() {
	pool := workerpool.New(10)

	startExecute := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local)
	endExecute := time.Now()
	duration := int(math.Ceil(endExecute.Sub(startExecute).Hours() / 24))

	initialDate := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local)
	for i := 0; i < duration; i++ {
		start := initialDate
		end := start.Add(24 * time.Hour)
		pool.Submit(h.run(&start, &end))
		initialDate = end
	}
	pool.StopWait()
}

func (h DailyProductionHandler) run(start, end *time.Time) func() {
	return func() {
		logger := logger.NewLogger(
			&logger.LoggerOption{
				LogName:     "logs/daily_production.log",
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
			logger.Errorf("[%v]Failed to connect to elasticsearch", start.Format(constant.YEAR_MONTH_DAY))
			return
		}

		masterSiteRepo, err := repo.NewMasterSiteRepo()
		if err != nil {
			logger.Error(err)
			return
		}

		solarRepo := repo.NewSolarRepo(elastic)
		serv := service.NewDailyProductionService(solarRepo, masterSiteRepo, logger)
		if err := serv.Run(start, end); err != nil {
			logger.Error(err)
		}
	}
}
