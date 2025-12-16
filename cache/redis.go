package cache

import (
	"github.com/hanson/go-toolbox/utils"
	"github.com/redis/go-redis/v9"
	"hertz-starter-kit/config"
	"log"
)

var Rdb *redis.Client

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.Cfg.RedisAddr,
		Password: config.Cfg.RedisPassword, // no password set
		DB:       config.Cfg.RedisDb,       // use default DB
	})
	ResetRedisLogger()
}

func ResetRedisLogger() {
	l := &logger{
		log: log.New(utils.GetMultiWriter(utils.HOUR), "redis: ", log.LstdFlags|log.Lshortfile),
	}
	redis.SetLogger(l)
}
