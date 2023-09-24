package handler

import (
	"fmt"
	"math"
	"time"

	"github.com/hugebear-io/true-solar-production/infra"
	"github.com/hugebear-io/true-solar-production/logger"
	"github.com/hugebear-io/true-solar-production/repo"
	"github.com/hugebear-io/true-solar-production/service"
	"github.com/sourcegraph/conc"
)

type DailyProductionHandler struct{}

func NewDailyProductionHandler() *DailyProductionHandler {
	return &DailyProductionHandler{}
}

func (h DailyProductionHandler) Run() {
	startExecute := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local)
	endExecute := time.Now()
	duration := int(math.Ceil(endExecute.Sub(startExecute).Hours() / 24))
	wg := conc.NewWaitGroup()

	initialDate := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local)
	for i := 0; i < duration; i++ {
		start := initialDate
		end := start.Add(24 * time.Hour)
		wg.Go(h.run(&start, &end))
		initialDate = end
	}
	wg.Wait()
}

func (h DailyProductionHandler) run(start, end *time.Time) func() {
	return func() {
		elastic, err := infra.NewElasticsearch()
		if err != nil {
			fmt.Println(err)
		}

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

		solarRepo := repo.NewSolarRepo(elastic)
		serv := service.NewDailyProductionService(solarRepo, logger)
		if err := serv.Run(start, end); err != nil {
			fmt.Println(err)
		}
	}
}
