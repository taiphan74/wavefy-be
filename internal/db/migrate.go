package db

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"wavefy-be/internal/model"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Role{}, &model.User{}); err != nil {
		return err
	}
	if err := seedRoles(db); err != nil {
		return err
	}
	return nil
}

func seedRoles(db *gorm.DB) error {
	roles := []model.Role{
		{
			ID:          uuid.New(),
			Name:        "USER",
			Description: "Default user role",
		},
		{
			ID:          uuid.New(),
			Name:        "ADMIN",
			Description: "Administrator role",
		},
	}

	for _, role := range roles {
		var existing model.Role
		err := db.Where("name = ?", role.Name).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := db.Create(&role).Error; err != nil {
			return err
		}
	}

	return nil
}
