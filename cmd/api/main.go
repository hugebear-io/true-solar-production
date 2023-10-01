package main

import (
	"github.com/hugebear-io/true-solar-production/config"
	"github.com/hugebear-io/true-solar-production/infra"
	"github.com/hugebear-io/true-solar-production/model"
)

func init() {
	config.InitConfig()
}

func main() {
	data, err := infra.NewGormDB("data.db")
	if err != nil {
		panic(err)
	}

	db, err := infra.NewGormDB()
	if err != nil {
		panic(err)
	}

	olds := []model.SiteRegionMapping{}
	if err := db.Find(&olds).Error; err != nil {
		panic(err)
	}

	if err := data.Create(olds).Error; err != nil {
		panic(err)
	}
}
