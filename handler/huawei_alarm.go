package handler

import (
	"time"

	"github.com/gammazero/workerpool"
	"github.com/hugebear-io/true-solar-production/constant"
	"github.com/hugebear-io/true-solar-production/infra"
	"github.com/hugebear-io/true-solar-production/logger"
	"github.com/hugebear-io/true-solar-production/model"
	"github.com/hugebear-io/true-solar-production/repo"
	"github.com/hugebear-io/true-solar-production/service"
)

type HuaweiAlarmHandler struct {
	logger logger.Logger
}

func NewHuaweiAlarmHandler() *HuaweiAlarmHandler {
	logger := logger.NewLogger(
		&logger.LoggerOption{
			LogName:     constant.GetLogName(constant.HUAWEI_ALARM_LOG_NAME),
			LogSize:     1024,
			LogAge:      90,
			LogBackup:   1,
			LogCompress: false,
			LogLevel:    logger.LOG_LEVEL_DEBUG,
			SkipCaller:  1,
		},
	)

	return &HuaweiAlarmHandler{logger: logger}
}

func (h *HuaweiAlarmHandler) Run() {
	h.logger = logger.NewLogger(
		&logger.LoggerOption{
			LogName:     constant.GetLogName(constant.HUAWEI_ALARM_LOG_NAME),
			LogSize:     1024,
			LogAge:      90,
			LogBackup:   1,
			LogCompress: false,
			LogLevel:    logger.LOG_LEVEL_DEBUG,
			SkipCaller:  1,
		},
	)
	defer h.logger.Close()

	db, err := infra.NewGormDB()
	if err != nil {
		h.logger.Error(err)
		return
	}

	credentialRepo := repo.NewHuaweiCredentialRepo(db)
	credentials, err := credentialRepo.GetCredentialsByOwner(constant.TRUE_OWNER)
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

func (h *HuaweiAlarmHandler) run(credential *model.HuaweiCredential) func() {
	return func() {
		now := time.Now()
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

		serv := service.NewHuaweiAlarmService(snmpRepo, rdb, h.logger)
		if err := serv.Run(credential); err != nil {
			h.logger.Errorf("[%v]Failed to run service: %v", credential.Username, err)
			return
		}
		h.logger.Infof("[%v] Finished in %v", credential.Username, time.Since(now).String())
	}
}

func (h *HuaweiAlarmHandler) Mock() {
	h.logger = logger.NewLogger(
		&logger.LoggerOption{
			LogName:     constant.GetLogName(constant.HUAWEI_ALARM_LOG_NAME),
			LogSize:     1024,
			LogAge:      90,
			LogBackup:   1,
			LogCompress: false,
			LogLevel:    logger.LOG_LEVEL_DEBUG,
			SkipCaller:  1,
		},
	)
	defer h.logger.Close()

	db, err := infra.NewGormDB()
	if err != nil {
		h.logger.Error(err)
		return
	}

	credentialRepo := repo.NewHuaweiCredentialRepo(db)
	credentials, err := credentialRepo.GetCredentialsByOwner(constant.TRUE_OWNER)
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

func (h *HuaweiAlarmHandler) mock(credential *model.HuaweiCredential) func() {
	return func() {
		now := time.Now()
		snmpRepo := repo.NewMockSnmpRepo()
		defer snmpRepo.Close()

		rdb, err := infra.NewRedis()
		if err != nil {
			h.logger.Errorf("[%v]Failed to connect to redis", credential.Username)
			return
		}
		defer rdb.Close()

		serv := service.NewHuaweiAlarmService(snmpRepo, rdb, h.logger)
		if err := serv.Run(credential); err != nil {
			h.logger.Errorf("[%v]Failed to run service: %v", credential.Username, err)
			return
		}
		h.logger.Infof("[%v] Finished in %v", credential.Username, time.Since(now).String())
	}
}
