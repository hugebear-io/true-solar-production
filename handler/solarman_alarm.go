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

type SolarmanAlarmHandler struct {
	logger logger.Logger
}

func NewSolarmanAlarmHandler() *SolarmanAlarmHandler {
	logger := logger.NewLogger(
		&logger.LoggerOption{
			LogName:     constant.GetLogName(constant.SOLARMAN_ALARM_LOG_NAME),
			LogSize:     1024,
			LogAge:      90,
			LogBackup:   1,
			LogCompress: false,
			LogLevel:    logger.LOG_LEVEL_DEBUG,
			SkipCaller:  1,
		},
	)

	return &SolarmanAlarmHandler{logger: logger}
}

func (h *SolarmanAlarmHandler) Run() {
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

func (h *SolarmanAlarmHandler) run(credential *model.SolarmanCredential) func() {
	return func() {
		snmp, err := infra.NewSnmp()
		if err != nil {
			h.logger.Errorf("[%v]Failed to connect to snmp", credential.Username)
			return
		}

		snmpRepo := repo.NewSnmpRepo(snmp)
		defer snmpRepo.Close()

		rdb, err := infra.NewRedis()
		if err != nil {
			h.logger.Errorf("[%v]Failed to connect to redis", credential.Username)
			return
		}
		defer rdb.Close()

		serv := service.NewSolarmanAlarmService(snmpRepo, rdb, h.logger)
		if err := serv.Run(credential); err != nil {
			h.logger.Errorf("[%v]Failed to run service: %v", credential.Username, err)
		}
	}
}

func (h *SolarmanAlarmHandler) Mock() {
	defer h.logger.Close()
	credentialRepo := repo.NewMockSolarmanCredentialRepo()
	credentials, err := credentialRepo.GetCredentials()
	if err != nil {
		h.logger.Error(err)
		return
	}

	pool := workerpool.New(1)
	for _, credential := range credentials {
		clone := credential
		pool.Submit(h.mock(&clone))
	}
	pool.StopWait()
}

func (h *SolarmanAlarmHandler) mock(credential *model.SolarmanCredential) func() {
	return func() {
		snmpRepo := repo.NewMockSnmpRepo()
		defer snmpRepo.Close()

		rdb, err := infra.NewRedis()
		if err != nil {
			h.logger.Error(err)
			return
		}
		defer rdb.Close()

		serv := service.NewSolarmanAlarmService(snmpRepo, rdb, h.logger)
		if err := serv.Run(credential); err != nil {
			h.logger.Error(err)
		}
	}
}
