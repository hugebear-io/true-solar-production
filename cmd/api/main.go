package main

import (
	"github.com/hugebear-io/true-solar-production/inverter/huawei"
	"github.com/hugebear-io/true-solar-production/util"
)

func main() {
	client, err := huawei.NewHuaweiClient(&huawei.HuaweiCredential{
		Username: "trueapi",
		Password: "Trueapi12@",
	})

	if err != nil {
		panic(err)
	}

	res, err := client.GetPlantList()
	if err != nil {
		panic(err)
	}

	util.PrintJSON(res)
}
