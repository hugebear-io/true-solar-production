package main

import (
	"fmt"
	"time"

	"github.com/hugebear-io/true-solar-production/infra"
	"github.com/hugebear-io/true-solar-production/logger"
	"github.com/hugebear-io/true-solar-production/repo"
	"github.com/hugebear-io/true-solar-production/service"
)

func main() {
	start := time.Date(2023, time.February, 1, 0, 0, 0, 0, time.Local)
	end := start.Add(24 * time.Hour)

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
		fmt.Println(err)
	}

	masterSiteRepo, err := repo.NewMasterSiteRepo()
	if err != nil {
		logger.Error(err)
		return
	}

	solarRepo := repo.NewSolarRepo(elastic)
	serv := service.NewDailyProductionService(solarRepo, masterSiteRepo, logger)
	if err := serv.Run(&start, &end); err != nil {
		logger.Error(err)
	}
}
