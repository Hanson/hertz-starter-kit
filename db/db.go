package db

import (
	"context"
	"fmt"
	"github.com/hanson/go-toolbox/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"hertz-starter-kit/config"
	log2 "hertz-starter-kit/utils/log"
	"hertz-starter-kit/utils/log/mysql_log"
	"log"
	"time"
)

var Db *gorm.DB

func InitDb(ctx context.Context) error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Cfg.DbUsername, config.Cfg.DbPassword, config.Cfg.DbAddress, config.Cfg.DbName)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger: mysql_log.New(log.New(utils.GetMultiWriter(utils.HOUR), "", log.LstdFlags),
			logger2.Config{
				SlowThreshold:             time.Second,
				Colorful:                  true,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      false,
				LogLevel:                  logger2.Info,
			}),
		NamingStrategy: DefaultNaming{},
	})
	if err != nil {
		log2.Errorf(ctx, "err: %+v", err)
		return err
	}

	db, err := Db.DB()
	if err != nil {
		log.Printf("err: %+v", err)
		return err
	}

	db.SetMaxOpenConns(5000)
	db.SetMaxIdleConns(1000)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(time.Hour)

	//if config.Cfg.DbSeparation {
	//	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Cfg.DbUsername, config.Cfg.DbPassword, config.Cfg.DbAddress, config.Cfg.DbName)
	//	Db.Use(dbresolver.Register(dbresolver.Config{
	//		Sources:           []gorm.Dialector{mysql.Open(dsn)},
	//		Replicas:          []gorm.Dialector{mysql.Open(dsn)},
	//		TraceResolverMode: true,
	//	}))
	//	return nil
	//}

	return nil
}

func UpdateLogger(split int) {
	Db = Db.Session(&gorm.Session{Logger: mysql_log.New(log.New(utils.GetMultiWriter(split), "", log.LstdFlags),
		logger2.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			LogLevel:                  logger2.Info,
		})})
}

func NewInstance(m interface{}, ctx context.Context) *gorm.DB {
	return Db.Model(m).Where("deleted_at = 0").Session(&gorm.Session{Context: ctx})
}
