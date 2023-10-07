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
	conf := config.GetConfig().Huawei
	collector := handler.NewHuaweiCollectorHandler()
	alarm := handler.NewHuaweiAlarmHandler()

	cron := gocron.NewScheduler(time.Local)
	cron.Cron(conf.CollectorCrontab).Do(collector.Run)
	cron.Cron(conf.NightCollectorCrontab).Do(collector.Run)
	cron.Cron(conf.AlarmCrontab).Do(alarm.Run)
	cron.StartBlocking()
	alarm.Mock()
}
