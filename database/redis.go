package database

import (
	"fmt"
	"passport-v4/config"

	"github.com/go-redis/redis"
)

var CfgDB *redis.Client

func initRedis() (err error) {
	CfgDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	_, err = CfgDB.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
