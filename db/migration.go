package db

import (
	"context"
	"hertz-starter-kit/utils/log"
)

func AutoMigrate(ctx context.Context) error {
	m := &MyMigrator{}
	m.DB = Db
	m.Migrator.Migrator.Dialector = Db.Dialector
	err := m.AutoMigrate()
	if err != nil {
		log.Errorf(ctx, "err: %+v", err)
		return err
	}

	return nil
}
