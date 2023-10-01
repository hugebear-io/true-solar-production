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
	conf := config.GetConfig().Solarman

	solarmanCollector := handler.NewSolarmanCollectorHandler()
	solarmanAlarm := handler.NewSolarmanAlarmHandler()

	cron := gocron.NewScheduler(time.Local)
	cron.Cron(conf.CollectorCrontab).Do(solarmanCollector.Run)
	cron.Cron(conf.NightCollectorCrontab).Do(solarmanCollector.Run)
	cron.Cron(conf.AlarmCrontab).Do(solarmanAlarm.Run)
	cron.StartBlocking()
}
