package handler

import (
	"github.com/gammazero/workerpool"
	"github.com/hugebear-io/true-solar-production/infra"
	"github.com/hugebear-io/true-solar-production/logger"
	"github.com/hugebear-io/true-solar-production/model"
	"github.com/hugebear-io/true-solar-production/repo"
	"github.com/hugebear-io/true-solar-production/service"
)

type SolarmanCollectorHandler struct {
}

func NewSolarmanCollectorHandler() *SolarmanCollectorHandler {
	return &SolarmanCollectorHandler{}
}

func (h SolarmanCollectorHandler) Mock() {
	logger := logger.NewLogger(
		&logger.LoggerOption{
			LogName:     "logs/solarman_collector.log",
			LogSize:     1024,
			LogAge:      90,
			LogBackup:   1,
			LogCompress: false,
			LogLevel:    logger.LOG_LEVEL_DEBUG,
			SkipCaller:  1,
		},
	)
	defer logger.Close()

	credentialRepo := repo.NewMockSolarmanCredentialRepo()
	credentials, err := credentialRepo.GetCredentials()
	if err != nil {
		logger.Error(err)
		return
	}

	pool := workerpool.New(len(credentials))
	for _, credential := range credentials {
		clone := credential
		pool.Submit(h.mock(&clone))
	}
	pool.StopWait()
}

func (h SolarmanCollectorHandler) mock(credential *model.SolarmanCredential) func() {
	return func() {
		logger := logger.NewLogger(
			&logger.LoggerOption{
				LogName:     "logs/solarman_collector.log",
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
			logger.Errorf("[%v]Failed to connect to elasticsearch", credential.Username)
			return
		}
		solarRepo := repo.NewSolarRepo(elastic)

		db, err := infra.NewGormDB()
		if err != nil {
			logger.Errorf("[%v]Failed to connect to gorm", credential.Username)
			return
		}
		siteRegionRepo := repo.NewSiteRegionMappingRepo(db)

		serv, err := service.NewSolarmanCollectorService(solarRepo, siteRegionRepo, logger)
		if err != nil {
			logger.Errorf("[%v]Failed to create service", credential.Username)
			return
		}

		if err := serv.Run(credential); err != nil {
			logger.Error(err)
		}
	}
}
