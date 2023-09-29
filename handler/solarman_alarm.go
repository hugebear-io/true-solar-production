package handler

import (
	"github.com/gammazero/workerpool"
	"github.com/hugebear-io/true-solar-production/infra"
	"github.com/hugebear-io/true-solar-production/logger"
	"github.com/hugebear-io/true-solar-production/model"
	"github.com/hugebear-io/true-solar-production/repo"
	"github.com/hugebear-io/true-solar-production/service"
)

type SolarmanAlarmHandler struct {
}

func NewSolarmanAlarmHandler() *SolarmanAlarmHandler {
	return &SolarmanAlarmHandler{}
}

func (h SolarmanAlarmHandler) Mock() {
	logger := logger.NewLogger(
		&logger.LoggerOption{
			LogName:     "logs/solarman_alarm.log",
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

	pool := workerpool.New(1)
	for _, credential := range credentials {
		clone := credential
		pool.Submit(h.mock(&clone))
	}
	pool.StopWait()
}

func (h SolarmanAlarmHandler) mock(credential *model.SolarmanCredential) func() {
	return func() {
		logger := logger.NewLogger(
			&logger.LoggerOption{
				LogName:     "logs/solarman_alarm.log",
				LogSize:     1024,
				LogAge:      90,
				LogBackup:   1,
				LogCompress: false,
				LogLevel:    logger.LOG_LEVEL_DEBUG,
				SkipCaller:  1,
			},
		)
		defer logger.Close()

		snmpRepo := repo.NewMockSnmpRepo()
		defer snmpRepo.Close()

		rdb, err := infra.NewRedis()
		if err != nil {
			logger.Error(err)
			return
		}
		defer rdb.Close()

		serv := service.NewSolarmanAlarmService(snmpRepo, rdb, logger)
		if err := serv.Run(credential); err != nil {
			logger.Error(err)
		}
	}
}
