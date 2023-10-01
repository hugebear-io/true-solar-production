package main

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/hugebear-io/true-solar-production/config"
	"github.com/hugebear-io/true-solar-production/handler"
	"github.com/hugebear-io/true-solar-production/util"
)

func init() {
	config.InitConfig()
}

func init() {
	util.SetTimezone()
}

func main() {
	lowPerformanceAlarmConf := config.GetConfig().LowPerformanceAlarm
	sumPerformanceAlarmConf := config.GetConfig().SumPerformanceAlarm

	lowPerformanceAlarmHdl := handler.NewLowPerformanceAlarmHandler()
	sumPerformanceAlarmHdl := handler.NewSumPerformanceAlarmHandler()

	cron := gocron.NewScheduler(time.Local)

	cron.Cron(lowPerformanceAlarmConf.Crontab).Do(lowPerformanceAlarmHdl.Mock)
	cron.Cron(sumPerformanceAlarmConf.Crontab).Do(sumPerformanceAlarmHdl.Mock)

	cron.StartBlocking()
}
