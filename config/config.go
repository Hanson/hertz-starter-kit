package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HostPort string

	AdminUsername string
	AdminPassword string

	AppUrl string

	// Db 配置
	DbUsername string
	DbPassword string
	DbAddress  string
	DbName     string

	AllowOrigins []string

	RedisAddr     string
	RedisPassword string
	RedisDb       int
}

var Cfg *Config

func LoadConfig() {
	Cfg = &Config{}

	err := godotenv.Load()
	if err == nil {
		Cfg.HostPort = os.Getenv("HOST_PORT")

		Cfg.AdminUsername = os.Getenv("ADMIN_USERNAME")
		Cfg.AdminPassword = os.Getenv("ADMIN_PASSWORD")

		Cfg.DbUsername = os.Getenv("DB_USERNAME")
		Cfg.DbPassword = os.Getenv("DB_PASSWORD")
		Cfg.DbAddress = os.Getenv("DB_ADDRESS")
		Cfg.DbName = os.Getenv("DB_NAME")

		Cfg.AppUrl = os.Getenv("APP_URL")

		Cfg.AllowOrigins = strings.Split(os.Getenv("ALLOW_ORIGINS"), ",")

		Cfg.RedisAddr = os.Getenv("REDIS_ADDR")
		Cfg.RedisPassword = os.Getenv("REDIS_PASSWORD")
		redisDb := os.Getenv("REDIS_DB")
		if redisDb != "" {
			Cfg.RedisDb, _ = strconv.Atoi(redisDb)
		} else {
			Cfg.RedisDb = 0
		}
	}
}
