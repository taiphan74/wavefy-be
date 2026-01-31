package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"wavefy-be/config"
	"wavefy-be/internal/mail"
	"wavefy-be/internal/model"
	"wavefy-be/internal/repository"
	"wavefy-be/internal/token"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidResetToken  = errors.New("invalid reset token")
	ErrMailNotConfigured  = errors.New("mail not configured")
)

type AuthService interface {
	Register(ctx context.Context, input CreateUserInput) (*model.User, *AuthToken, error)
	Login(ctx context.Context, input LoginInput) (*model.User, *AuthToken, error)
	Refresh(ctx context.Context, refreshToken string) (*model.User, *AuthToken, error)
	Logout(ctx context.Context, refreshToken string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, password string) error
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
	resetStore   token.PasswordResetTokenStore
	mailer       *mail.Service
	cfg          config.AuthConfig
}

func NewAuthService(userService UserService, userRepo repository.UserRepository, roleRepo repository.RoleRepository, refreshStore token.RefreshTokenStore, resetStore token.PasswordResetTokenStore, mailer *mail.Service, cfg config.AuthConfig) AuthService {
	return &authService{
		userService:  userService,
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		refreshStore: refreshStore,
		resetStore:   resetStore,
		mailer:       mailer,
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

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return ErrInvalidInput
	}
	if s.resetStore == nil || s.mailer == nil {
		return ErrMailNotConfigured
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	resetToken, err := s.resetStore.Create(ctx, user.ID.String())
	if err != nil {
		return err
	}

	subject := "Reset your password"
	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", resetToken)
	textBody := "Reset your password using this link: " + resetURL
	htmlBody := fmt.Sprintf(`<p>Reset your password: <a href="%s">Reset Password</a></p>`, resetURL)
	if err := s.mailer.Send(user.Email, subject, textBody, htmlBody); err != nil {
		_ = s.resetStore.Revoke(ctx, resetToken)
		return err
	}
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, resetToken, password string) error {
	if strings.TrimSpace(resetToken) == "" || password == "" {
		return ErrInvalidInput
	}
	if s.resetStore == nil {
		return ErrInvalidResetToken
	}

	userID, err := s.resetStore.Verify(ctx, resetToken)
	if err != nil {
		return ErrInvalidResetToken
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return ErrInvalidResetToken
	}

	user, err := s.userRepo.GetByID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidResetToken
		}
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	_ = s.resetStore.Revoke(ctx, resetToken)
	return nil
}
