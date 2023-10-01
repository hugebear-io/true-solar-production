package handler

import (
	"github.com/hugebear-io/true-solar-production/constant"
	"github.com/hugebear-io/true-solar-production/infra"
	"github.com/hugebear-io/true-solar-production/logger"
	"github.com/hugebear-io/true-solar-production/repo"
	"github.com/hugebear-io/true-solar-production/service"
)

type LowPerformanceAlarmHandler struct {
	logger logger.Logger
}

func NewLowPerformanceAlarmHandler() *LowPerformanceAlarmHandler {
	logger := logger.NewLogger(
		&logger.LoggerOption{
			LogName:     constant.GetLogName(constant.LOW_PERFORMANCE_ALARM_LOG_NAME),
			LogSize:     1024,
			LogAge:      90,
			LogBackup:   1,
			LogCompress: false,
			LogLevel:    logger.LOG_LEVEL_DEBUG,
			SkipCaller:  1,
		},
	)

	return &LowPerformanceAlarmHandler{
		logger: logger,
	}
}

func (h *LowPerformanceAlarmHandler) Run() {
	snmp, err := infra.NewSnmp()
	if err != nil {
		h.logger.Errorf("Failed to create snmp: %v", err)
		return
	}

	elastic, err := infra.NewElasticsearch()
	if err != nil {
		h.logger.Errorf("Failed to create elasticsearch: %v", err)
		return
	}

	db, err := infra.NewGormDB()
	if err != nil {
		h.logger.Errorf("Failed to create gorm db: %v", err)
		return
	}

	solarRepo := repo.NewSolarRepo(elastic)
	installedCapacityRepo := repo.NewInstalledCapacityRepo(db)
	performanceAlarmConfigRepo := repo.NewPerformanceAlarmConfigRepo(db)
	snmpRepo := repo.NewSnmpRepo(snmp)
	defer snmpRepo.Close()
	defer h.logger.Close()

	serv := service.NewLowPerformanceAlarmService(solarRepo, installedCapacityRepo, performanceAlarmConfigRepo, snmpRepo, h.logger)

	h.logger.Info("Running low performance alarm service")
	if err := serv.Run(); err != nil {
		h.logger.Errorf("Failed to run low performance alarm service: %v", err)
		return
	}
}

func (h *LowPerformanceAlarmHandler) Mock() {
	elastic, err := infra.NewElasticsearch()
	if err != nil {
		h.logger.Errorf("Failed to create elasticsearch: %v", err)
		return
	}

	db, err := infra.NewGormDB()
	if err != nil {
		h.logger.Errorf("Failed to create gorm db: %v", err)
		return
	}

	solarRepo := repo.NewSolarRepo(elastic)
	installedCapacityRepo := repo.NewInstalledCapacityRepo(db)
	performanceAlarmConfigRepo := repo.NewPerformanceAlarmConfigRepo(db)
	snmpRepo := repo.NewMockSnmpRepo()
	defer snmpRepo.Close()
	defer h.logger.Close()

	serv := service.NewLowPerformanceAlarmService(solarRepo, installedCapacityRepo, performanceAlarmConfigRepo, snmpRepo, h.logger)

	h.logger.Info("Running low performance alarm service")
	if err := serv.Run(); err != nil {
		h.logger.Errorf("Failed to run low performance alarm service: %v", err)
		return
	}
	h.logger.Info("Finished low performance alarm service")
}
