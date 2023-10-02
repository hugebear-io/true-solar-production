package handler

import (
	"github.com/gammazero/workerpool"
	"github.com/hugebear-io/true-solar-production/constant"
	"github.com/hugebear-io/true-solar-production/infra"
	"github.com/hugebear-io/true-solar-production/logger"
	"github.com/hugebear-io/true-solar-production/model"
	"github.com/hugebear-io/true-solar-production/repo"
	"github.com/hugebear-io/true-solar-production/service"
)

type SolarmanCollectorHandler struct {
	logger logger.Logger
}

func NewSolarmanCollectorHandler() *SolarmanCollectorHandler {
	logger := logger.NewLogger(
		&logger.LoggerOption{
			LogName:     constant.GetLogName(constant.SOLARMAN_COLLECTOR_LOG_NAME),
			LogSize:     1024,
			LogAge:      90,
			LogBackup:   1,
			LogCompress: false,
			LogLevel:    logger.LOG_LEVEL_DEBUG,
			SkipCaller:  1,
		},
	)

	return &SolarmanCollectorHandler{logger: logger}
}

func (h *SolarmanCollectorHandler) Run() {
	defer h.logger.Close()
	db, err := infra.NewGormDB()
	if err != nil {
		h.logger.Error(err)
		return
	}

	credentialRepo := repo.NewSolarmanCredentialRepo(db)
	credentials, err := credentialRepo.GetCredentials()
	if err != nil {
		h.logger.Error(err)
		return
	}

	pool := workerpool.New(len(credentials))
	for _, credential := range credentials {
		clone := credential
		pool.Submit(h.run(&clone))
	}
	pool.StopWait()
}

func (h *SolarmanCollectorHandler) run(credential *model.SolarmanCredential) func() {
	return func() {
		elastic, err := infra.NewElasticsearch()
		if err != nil {
			h.logger.Errorf("[%v]Failed to connect to elasticsearch", credential.Username)
			return
		}
		solarRepo := repo.NewSolarRepo(elastic)

		db, err := infra.NewGormDB()
		if err != nil {
			h.logger.Errorf("[%v]Failed to connect to gorm", credential.Username)
			return
		}
		siteRegionRepo := repo.NewSiteRegionMappingRepo(db)

		serv, err := service.NewSolarmanCollectorService(solarRepo, siteRegionRepo, h.logger)
		if err != nil {
			h.logger.Errorf("[%v]Failed to create service", credential.Username)
			return
		}

		if err := serv.Run(credential); err != nil {
			h.logger.Errorf("[%v]Failed to run service: %v", credential.Username, err)
			return
		}

		h.logger.Infof("[%v]Finished", credential.Username)
	}
}

func (h *SolarmanCollectorHandler) Mock() {
	defer h.logger.Close()
	credentialRepo := repo.NewMockSolarmanCredentialRepo()
	credentials, err := credentialRepo.GetCredentials()
	if err != nil {
		h.logger.Error(err)
		return
	}

	pool := workerpool.New(len(credentials))
	for _, credential := range credentials {
		clone := credential
		pool.Submit(h.mock(&clone))
	}
	pool.StopWait()
}

func (h *SolarmanCollectorHandler) mock(credential *model.SolarmanCredential) func() {
	return func() {
		elastic, err := infra.NewElasticsearch()
		if err != nil {
			h.logger.Errorf("[%v]Failed to connect to elasticsearch", credential.Username)
			return
		}
		solarRepo := repo.NewSolarRepo(elastic)

		db, err := infra.NewGormDB()
		if err != nil {
			h.logger.Errorf("[%v]Failed to connect to gorm", credential.Username)
			return
		}
		siteRegionRepo := repo.NewSiteRegionMappingRepo(db)

		serv, err := service.NewSolarmanCollectorService(solarRepo, siteRegionRepo, h.logger)
		if err != nil {
			h.logger.Errorf("[%v]Failed to create service", credential.Username)
			return
		}

		if err := serv.Run(credential); err != nil {
			h.logger.Error(err)
		}
	}
}
