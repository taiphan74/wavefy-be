package model

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"size:50;uniqueIndex;not null"`
	Description string    `gorm:"size:255"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
