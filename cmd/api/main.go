package main

import (
	"fmt"

	"github.com/hugebear-io/true-solar-production/infra"
	"github.com/hugebear-io/true-solar-production/model"
	"github.com/hugebear-io/true-solar-production/util"
)

type Schema struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

func main() {
	legacy, err := infra.NewGormDB("legacy.db")
	if err != nil {
		panic(err)
	}

	schema := []Schema{}
	statement := "SELECT username, password, app_id, app_secret FROM data_collector_config WHERE vendor_type='invt';"
	if err := legacy.Raw(statement).Scan(&schema).Error; err != nil {
		panic(err)
	}
	util.PrintJSON(map[string]interface{}{"result": schema})

	credentials := make([]model.SolarmanCredential, 0)
	for _, item := range schema {
		credentials = append(credentials, model.SolarmanCredential{
			Username:  item.Username,
			Password:  item.Password,
			AppID:     item.AppID,
			AppSecret: item.AppSecret,
		})
	}

	db, err := infra.NewGormDB()
	if err != nil {
		panic(err)
	}

	if err := db.Create(&credentials).Error; err != nil {
		panic(err)
	}

	fmt.Println("done")
}
