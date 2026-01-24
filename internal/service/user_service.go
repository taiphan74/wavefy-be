package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"wavefy-be/internal/model"
	"wavefy-be/internal/repository"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrEmailExists  = errors.New("email already exists")
	ErrNotFound     = errors.New("not found")
)

type CreateUserInput struct {
	Email    string
	Password string
}

type UserService interface {
	Create(ctx context.Context, input CreateUserInput) (*model.User, error)
	Get(ctx context.Context, id uuid.UUID) (*model.User, error)
	List(ctx context.Context, limit, offset int) ([]model.User, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	repo     repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewUserService(repo repository.UserRepository, roleRepo repository.RoleRepository) UserService {
	return &userService{repo: repo, roleRepo: roleRepo}
}

func (s *userService) Create(ctx context.Context, input CreateUserInput) (*model.User, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	if input.Email == "" || input.Password == "" {
		return nil, ErrInvalidInput
	}

	if _, err := s.repo.GetByEmail(ctx, input.Email); err == nil {
		return nil, ErrEmailExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: string(hash),
	}

	role, err := s.roleRepo.GetByName(ctx, "USER")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidInput
		}
		return nil, err
	}
	user.RoleID = role.ID

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

type UpdateUserInput struct {
	FirstName *string
	LastName  *string
	Email     *string
	Password  *string
}

func (s *userService) Get(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) List(ctx context.Context, limit, offset int) ([]model.User, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if input.FirstName != nil {
		user.FirstName = strings.TrimSpace(*input.FirstName)
	}
	if input.LastName != nil {
		user.LastName = strings.TrimSpace(*input.LastName)
	}
	if input.Email != nil {
		email := strings.TrimSpace(strings.ToLower(*input.Email))
		if email == "" {
			return nil, ErrInvalidInput
		}
		if existing, err := s.repo.GetByEmail(ctx, email); err == nil && existing.ID != user.ID {
			return nil, ErrEmailExists
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		user.Email = email
	}
	if input.Password != nil {
		if *input.Password == "" {
			return nil, ErrInvalidInput
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(hash)
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}
