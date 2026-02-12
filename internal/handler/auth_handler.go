package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"wavefy-be/config"
	"wavefy-be/helper"
	"wavefy-be/internal/dto"
	"wavefy-be/internal/service"
)

type AuthHandler struct {
	service service.AuthService
	cfg     config.AuthConfig
}

const refreshCookieName = "refresh_token"
const refreshCookiePath = "/api/auth/refresh"

func NewAuthHandler(service service.AuthService, cfg config.AuthConfig) *AuthHandler {
	return &AuthHandler{
		service: service,
		cfg:     cfg,
	}
}

// Register godoc
// @Summary      Register
// @Description  Register a new account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Register"
// @Success      200 {object} helper.Response{data=dto.AuthResponse}
// @Failure      400 {object} helper.Response
// @Failure      409 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, token, err := h.service.Register(c.Request.Context(), service.CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			helper.RespondError(c, http.StatusBadRequest, err.Error())
		case service.ErrEmailExists:
			helper.RespondError(c, http.StatusConflict, err.Error())
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	h.setRefreshCookie(c, token.RefreshToken)

	helper.RespondOK(c, dto.AuthResponse{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresAt:   token.ExpiresAt.Format(time.RFC3339),
		User:        mapUserResponse(user),
	})
}

// Login godoc
// @Summary      Login
// @Description  Login with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login"
// @Success      200 {object} helper.Response{data=dto.AuthResponse}
// @Failure      400 {object} helper.Response
// @Failure      401 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, token, err := h.service.Login(c.Request.Context(), service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			helper.RespondError(c, http.StatusUnauthorized, err.Error())
		case service.ErrTooManyAttempts:
			helper.RespondError(c, http.StatusTooManyRequests, err.Error())
		case service.ErrEmailNotVerified:
			helper.RespondError(c, http.StatusForbidden, err.Error())
		case service.ErrMailNotConfigured:
			helper.RespondError(c, http.StatusServiceUnavailable, err.Error())
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	h.setRefreshCookie(c, token.RefreshToken)

	helper.RespondOK(c, dto.AuthResponse{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresAt:   token.ExpiresAt.Format(time.RFC3339),
		User:        mapUserResponse(user),
	})
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Rotate refresh token and issue a new access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshRequest true "Refresh"
// @Success      200 {object} helper.Response{data=dto.AuthResponse}
// @Failure      400 {object} helper.Response
// @Failure      401 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshCookieName)
	if err != nil || refreshToken == "" {
		helper.RespondError(c, http.StatusUnauthorized, "missing refresh token")
		return
	}

	user, token, err := h.service.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			helper.RespondError(c, http.StatusUnauthorized, err.Error())
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	h.setRefreshCookie(c, token.RefreshToken)

	helper.RespondOK(c, dto.AuthResponse{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresAt:   token.ExpiresAt.Format(time.RFC3339),
		User:        mapUserResponse(user),
	})
}

// ForgotPassword godoc
// @Summary      Forgot password
// @Description  Send reset password email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.ForgotPasswordRequest true "Forgot password"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      503 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.ForgotPassword(c.Request.Context(), req.Email); err != nil {
		switch err {
		case service.ErrInvalidInput:
			helper.RespondError(c, http.StatusBadRequest, err.Error())
		case service.ErrMailNotConfigured:
			helper.RespondError(c, http.StatusServiceUnavailable, err.Error())
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondOK(c, gin.H{"sent": true})
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Reset password by token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.ResetPasswordRequest true "Reset password"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      401 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), req.Token, req.Password); err != nil {
		switch err {
		case service.ErrInvalidInput:
			helper.RespondError(c, http.StatusBadRequest, err.Error())
		case service.ErrInvalidResetToken:
			helper.RespondError(c, http.StatusUnauthorized, err.Error())
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondOK(c, gin.H{"reset": true})
}

// VerifyEmail godoc
// @Summary      Verify email
// @Description  Verify email by token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.VerifyEmailRequest true "Verify email"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      401 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.VerifyEmail(c.Request.Context(), req.Token); err != nil {
		switch err {
		case service.ErrInvalidVerifyToken:
			helper.RespondError(c, http.StatusUnauthorized, err.Error())
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondOK(c, gin.H{"verified": true})
}

// Logout godoc
// @Summary      Logout
// @Description  Revoke refresh token and clear refresh cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshRequest false "Refresh"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      401 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshCookieName)
	if err != nil || refreshToken == "" {
		helper.RespondError(c, http.StatusUnauthorized, "missing refresh token")
		return
	}

	if err := h.service.Logout(c.Request.Context(), refreshToken); err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			helper.RespondError(c, http.StatusUnauthorized, err.Error())
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	h.clearRefreshCookie(c)
	helper.RespondOK(c, gin.H{"logged_out": true})
}

func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string) {
	if token == "" || h.cfg.RefreshTokenTTL <= 0 {
		return
	}

	expiresAt := time.Now().UTC().Add(h.cfg.RefreshTokenTTL)
	secure := isSecureRequest(c)

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     refreshCookiePath,
		MaxAge:   int(h.cfg.RefreshTokenTTL.Seconds()),
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {
	secure := isSecureRequest(c)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     refreshCookiePath,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0).UTC(),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func isSecureRequest(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
	}
	if strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https") {
		return true
	}
	return false
}
