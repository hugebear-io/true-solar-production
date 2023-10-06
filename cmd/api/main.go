package main

import (
	"fmt"

	"github.com/hugebear-io/true-solar-production/config"
	"github.com/hugebear-io/true-solar-production/constant"
	"github.com/hugebear-io/true-solar-production/infra"
	"github.com/hugebear-io/true-solar-production/repo"
)

func init() {
	config.InitConfig()
}

func init() {
}

func main() {
	db, err := infra.NewGormDB()
	if err != nil {
		panic(err)
	}

	repo := repo.NewHuaweiCredentialRepo(db)
	data, err := repo.GetCredentialsByOwner(constant.TRUE_OWNER)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)
}
