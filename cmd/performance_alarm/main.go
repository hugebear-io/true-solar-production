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
	// conf := config.GetConfig().DailyPerformanceAlarm
	dailyPerformanceAlarm := handler.NewDailyPerformanceAlarmHandler()

	cron := gocron.NewScheduler(time.Local)
	cron.Every(1).Minute().Do(dailyPerformanceAlarm.Mock)
	cron.StartBlocking()
}
