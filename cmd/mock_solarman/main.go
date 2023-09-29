package main

import (
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
	hdl := handler.NewSolarmanAlarmHandler()
	hdl.Mock()
}
