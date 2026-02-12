package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"wavefy-be/internal/model"
	"wavefy-be/internal/repository"
)

type CreateTrackInput struct {
	ArtistUserID string
	AlbumID      *string
	Title        string
	AudioURL     string
	ImageURL     *string
	DurationSec  int
	IsPublic     *bool
}

type UpdateTrackInput struct {
	AlbumID     *string
	Title       *string
	AudioURL    *string
	ImageURL    *string
	DurationSec *int
	IsPublic    *bool
}

type TrackService interface {
	Create(ctx context.Context, input CreateTrackInput) (*model.Track, error)
	Get(ctx context.Context, id uuid.UUID) (*model.Track, error)
	List(ctx context.Context, limit, offset int) ([]model.Track, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateTrackInput) (*model.Track, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type trackService struct {
	repo     repository.TrackRepository
	userRepo repository.UserRepository
}

func NewTrackService(repo repository.TrackRepository, userRepo repository.UserRepository) TrackService {
	return &trackService{repo: repo, userRepo: userRepo}
}

func (s *trackService) Create(ctx context.Context, input CreateTrackInput) (*model.Track, error) {
	artistID, err := uuid.Parse(strings.TrimSpace(input.ArtistUserID))
	if err != nil {
		return nil, ErrInvalidInput
	}

	if _, err := s.userRepo.GetByID(ctx, artistID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidInput
		}
		return nil, err
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, ErrInvalidInput
	}

	audioURL := strings.TrimSpace(input.AudioURL)
	if audioURL == "" {
		return nil, ErrInvalidInput
	}

	var imageURL *string
	if input.ImageURL != nil {
		value := strings.TrimSpace(*input.ImageURL)
		if value == "" {
			return nil, ErrInvalidInput
		}
		imageURL = &value
	}

	if input.DurationSec <= 0 {
		return nil, ErrInvalidInput
	}

	var albumID *uuid.UUID
	if input.AlbumID != nil {
		value := strings.TrimSpace(*input.AlbumID)
		if value != "" {
			parsed, err := uuid.Parse(value)
			if err != nil {
				return nil, ErrInvalidInput
			}
			albumID = &parsed
		}
	}

	isPublic := true
	if input.IsPublic != nil {
		isPublic = *input.IsPublic
	}

	track := &model.Track{
		ID:           uuid.New(),
		ArtistUserID: artistID,
		AlbumID:      albumID,
		Title:        title,
		AudioURL:     audioURL,
		ImageURL:     imageURL,
		DurationSec:  input.DurationSec,
		IsPublic:     isPublic,
		PlayCount:    0,
	}

	if err := s.repo.Create(ctx, track); err != nil {
		return nil, err
	}

	return track, nil
}

func (s *trackService) Get(ctx context.Context, id uuid.UUID) (*model.Track, error) {
	track, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return track, nil
}

func (s *trackService) List(ctx context.Context, limit, offset int) ([]model.Track, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *trackService) Update(ctx context.Context, id uuid.UUID, input UpdateTrackInput) (*model.Track, error) {
	track, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return nil, ErrInvalidInput
		}
		track.Title = title
	}

	if input.AudioURL != nil {
		audioURL := strings.TrimSpace(*input.AudioURL)
		if audioURL == "" {
			return nil, ErrInvalidInput
		}
		track.AudioURL = audioURL
	}

	if input.ImageURL != nil {
		value := strings.TrimSpace(*input.ImageURL)
		if value == "" {
			track.ImageURL = nil
		} else {
			track.ImageURL = &value
		}
	}

	if input.DurationSec != nil {
		if *input.DurationSec <= 0 {
			return nil, ErrInvalidInput
		}
		track.DurationSec = *input.DurationSec
	}

	if input.IsPublic != nil {
		track.IsPublic = *input.IsPublic
	}

	if input.AlbumID != nil {
		value := strings.TrimSpace(*input.AlbumID)
		if value == "" {
			track.AlbumID = nil
		} else {
			parsed, err := uuid.Parse(value)
			if err != nil {
				return nil, ErrInvalidInput
			}
			track.AlbumID = &parsed
		}
	}

	if err := s.repo.Update(ctx, track); err != nil {
		return nil, err
	}

	return track, nil
}

func (s *trackService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}
