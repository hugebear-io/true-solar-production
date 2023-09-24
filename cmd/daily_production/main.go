package main

import (
	"time"

	"github.com/hugebear-io/true-solar-production/config"
	"github.com/hugebear-io/true-solar-production/handler"
)

func init() {
	config.InitConfig()
}

func init() {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	time.Local = loc
}

func main() {
	dailyProduction := handler.NewDailyProductionHandler()
	dailyProduction.Run()
}
