package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	FirstName    string    `gorm:"size:100"`
	LastName     string    `gorm:"size:100"`
	Email        string    `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	IsActive     bool      `gorm:"default:false"`
	RoleID       uuid.UUID `gorm:"type:uuid;not null;index:idx_users_role_id"`
	Role         Role      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
