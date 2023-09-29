package main

import (
	"context"

	"github.com/hugebear-io/true-solar-production/infra"
)

func main() {
	rdb, _ := infra.NewRedis()
	ctx := context.TODO()
	if err := rdb.FlushAll(ctx).Err(); err != nil {
		panic(err)
	}

}
