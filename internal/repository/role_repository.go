package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"wavefy-be/internal/model"
)

type RoleRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Role, error)
	GetByName(ctx context.Context, name string) (*model.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
