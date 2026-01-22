package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"wavefy-be/config"
	"wavefy-be/internal/model"
	"wavefy-be/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService interface {
	Register(ctx context.Context, input CreateUserInput) (*model.User, *AuthToken, error)
	Login(ctx context.Context, input LoginInput) (*model.User, *AuthToken, error)
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthToken struct {
	AccessToken string
	ExpiresAt   time.Time
	TokenType   string
}

type authService struct {
	userService UserService
	userRepo    repository.UserRepository
	cfg         config.AuthConfig
}

func NewAuthService(userService UserService, userRepo repository.UserRepository, cfg config.AuthConfig) AuthService {
	return &authService{
		userService: userService,
		userRepo:    userRepo,
		cfg:         cfg,
	}
}

func (s *authService) Register(ctx context.Context, input CreateUserInput) (*model.User, *AuthToken, error) {
	user, err := s.userService.Create(ctx, input)
	if err != nil {
		return nil, nil, err
	}
	token, err := s.issueToken(user.ID.String())
	if err != nil {
		return nil, nil, err
	}
	return user, token, nil
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

	token, err := s.issueToken(user.ID.String())
	if err != nil {
		return nil, nil, err
	}
	return user, token, nil
}

func (s *authService) issueToken(subject string) (*AuthToken, error) {
	expiresAt := time.Now().UTC().Add(s.cfg.AccessTokenTTL)
	claims := jwt.RegisteredClaims{
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Issuer:    s.cfg.AccessTokenIss,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &AuthToken{
		AccessToken: signed,
		ExpiresAt:   expiresAt,
		TokenType:   "Bearer",
	}, nil
}
