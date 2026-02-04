package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Track struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey"`
	ArtistUserID uuid.UUID  `gorm:"type:uuid;not null;index:idx_tracks_artist_user_id"`
	ArtistUser   User       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	AlbumID      *uuid.UUID `gorm:"type:uuid;index:idx_tracks_album_id"`
	Title        string     `gorm:"size:255;not null"`
	AudioURL     string     `gorm:"size:800;not null"`
	DurationSec  int        `gorm:"not null"`
	IsPublic     bool       `gorm:"not null;default:true"`
	PlayCount    int64      `gorm:"type:bigint;not null;default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
