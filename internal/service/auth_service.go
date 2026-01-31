package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"wavefy-be/config"
	"wavefy-be/internal/model"
	"wavefy-be/internal/repository"
	"wavefy-be/internal/token"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService interface {
	Register(ctx context.Context, input CreateUserInput) (*model.User, *AuthToken, error)
	Login(ctx context.Context, input LoginInput) (*model.User, *AuthToken, error)
	Refresh(ctx context.Context, refreshToken string) (*model.User, *AuthToken, error)
	Logout(ctx context.Context, refreshToken string) error
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	TokenType    string
}

type authService struct {
	userService  UserService
	userRepo     repository.UserRepository
	roleRepo     repository.RoleRepository
	refreshStore token.RefreshTokenStore
	cfg          config.AuthConfig
}

func NewAuthService(userService UserService, userRepo repository.UserRepository, roleRepo repository.RoleRepository, refreshStore token.RefreshTokenStore, cfg config.AuthConfig) AuthService {
	return &authService{
		userService:  userService,
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		refreshStore: refreshStore,
		cfg:          cfg,
	}
}

func (s *authService) Register(ctx context.Context, input CreateUserInput) (*model.User, *AuthToken, error) {
	user, err := s.userService.Create(ctx, input)
	if err != nil {
		return nil, nil, err
	}
	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil {
		return nil, nil, err
	}
	accessToken, expiresAt, err := token.IssueAccessToken(s.cfg, user.ID.String(), role.Name)
	if err != nil {
		return nil, nil, err
	}
	authToken := &AuthToken{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
		TokenType:   "Bearer",
	}
	refreshToken, err := s.refreshStore.Create(ctx, user.ID.String())
	if err != nil {
		return nil, nil, err
	}
	authToken.RefreshToken = refreshToken
	return user, authToken, nil
}

func (s *authService) Login(ctx context.Context, input LoginInput) (*model.User, *AuthToken, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || input.Password == "" {
		return nil, nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil {
		return nil, nil, err
	}
	accessToken, expiresAt, err := token.IssueAccessToken(s.cfg, user.ID.String(), role.Name)
	if err != nil {
		return nil, nil, err
	}
	authToken := &AuthToken{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
		TokenType:   "Bearer",
	}
	refreshToken, err := s.refreshStore.Create(ctx, user.ID.String())
	if err != nil {
		return nil, nil, err
	}
	authToken.RefreshToken = refreshToken
	return user, authToken, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*model.User, *AuthToken, error) {
	userID, err := s.refreshStore.Verify(ctx, refreshToken)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.GetByID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, err
	}

	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil {
		return nil, nil, err
	}

	accessToken, expiresAt, err := token.IssueAccessToken(s.cfg, user.ID.String(), role.Name)
	if err != nil {
		return nil, nil, err
	}
	authToken := &AuthToken{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
		TokenType:   "Bearer",
	}

	newRefresh, err := s.refreshStore.Create(ctx, user.ID.String())
	if err != nil {
		return nil, nil, err
	}
	_ = s.refreshStore.Revoke(ctx, refreshToken)
	authToken.RefreshToken = newRefresh
	return user, authToken, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return ErrInvalidCredentials
	}
	return s.refreshStore.Revoke(ctx, refreshToken)
}
