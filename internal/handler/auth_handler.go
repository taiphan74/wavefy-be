package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"wavefy-be/helper"
	"wavefy-be/internal/dto"
	"wavefy-be/internal/service"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
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
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondOK(c, dto.AuthResponse{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresAt:   token.ExpiresAt.Format(time.RFC3339),
		User:        mapUserResponse(user),
	})
}
