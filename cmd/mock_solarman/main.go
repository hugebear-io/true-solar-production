package main

import (
	"github.com/hugebear-io/true-solar-production/config"
	"github.com/hugebear-io/true-solar-production/infra"
	"github.com/hugebear-io/true-solar-production/util"
)

func init() {
	config.InitConfig()
}

func init() {
	util.SetTimezone()
}

func main() {
	// hdl := handler.NewSolarmanAlarmHandler()
	// hdl.Mock()
	// hdl := handler.NewSolarmanCollectorHandler()
	// hdl.Mock()
	_, err := infra.NewRedis()
	if err != nil {
		panic(err)
	}
}
