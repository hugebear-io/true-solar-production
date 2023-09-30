package main

import (
	"github.com/hugebear-io/true-solar-production/repo"
	"github.com/hugebear-io/true-solar-production/util"
)

func main() {
	repo, err := repo.NewMasterSiteRepo()
	if err != nil {
		panic(err)
	}

	data := repo.ExportToMap()
	util.PrintJSON(map[string]interface{}{"x": data})
}
