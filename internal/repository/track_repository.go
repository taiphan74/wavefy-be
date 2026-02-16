package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"wavefy-be/internal/model"
)

type TrackRepository interface {
	Create(ctx context.Context, track *model.Track) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Track, error)
	List(ctx context.Context, limit, offset int) ([]model.Track, error)
	Update(ctx context.Context, track *model.Track) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type trackRepository struct {
	db *gorm.DB
}

func NewTrackRepository(db *gorm.DB) TrackRepository {
	return &trackRepository{db: db}
}

func (r *trackRepository) Create(ctx context.Context, track *model.Track) error {
	return r.db.WithContext(ctx).Create(track).Error
}

func (r *trackRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Track, error) {
	var track model.Track
	err := r.db.WithContext(ctx).Preload("ArtistUser").First(&track, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &track, nil
}

func (r *trackRepository) List(ctx context.Context, limit, offset int) ([]model.Track, error) {
	var tracks []model.Track
	err := r.db.WithContext(ctx).Preload("ArtistUser").Limit(limit).Offset(offset).Order("created_at desc").Find(&tracks).Error
	return tracks, err
}

func (r *trackRepository) Update(ctx context.Context, track *model.Track) error {
	return r.db.WithContext(ctx).Save(track).Error
}

func (r *trackRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Track{}, "id = ?", id).Error
}
