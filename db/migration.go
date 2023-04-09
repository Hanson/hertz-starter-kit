package db

import (
	"log"
)

func AutoMigrate() error {
	m := &MyMigrator{}
	m.DB = Db
	m.Migrator.Migrator.Dialector = Db.Dialector
	err := m.AutoMigrate()
	if err != nil {
		log.Printf("err: %+v", err)
		return err
	}

	return nil
}
