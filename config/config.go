package config

import (
	"github.com/joho/godotenv"
	"os"
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
	}
}
