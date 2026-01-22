package db

import (
	"gorm.io/gorm"

	"wavefy-be/internal/model"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{})
}
