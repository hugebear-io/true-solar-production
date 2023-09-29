package infra

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/hugebear-io/true-solar-production/config"
)

func NewRedis() (*redis.Client, error) {
	conf := config.GetConfig().Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + conf.Port,
		Username: conf.Username,
		Password: conf.Password,
		DB:       conf.DB,
	})

	ctx := context.TODO()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
