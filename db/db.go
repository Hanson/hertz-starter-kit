package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"hertz-starter-kit/config"
	"log"
	"os"
	"time"
)

var Db *gorm.DB

func InitDb() error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Cfg.DbUsername, config.Cfg.DbPassword, config.Cfg.DbAddress, config.Cfg.DbName)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: false,
		Logger: logger.New(log.New(os.Stdout, "", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				Colorful:                  true,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      false,
				LogLevel:                  logger.Info,
			}),
		NamingStrategy: DefaultNaming{},
	})
	if err != nil {
		log.Printf("err: %+v", err)
		return err
	}

	return nil
}

func NewInstance(m interface{}) *gorm.DB {
	return Db.Model(m).Where("deleted_at = 0")
}
