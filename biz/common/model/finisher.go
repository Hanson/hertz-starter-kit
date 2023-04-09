package model

import (
	"gorm.io/gorm"
	"time"
)

func SoftDelete(m *gorm.DB) *gorm.DB {
	return m.Update("deleted_at", time.Now().Unix())
}
