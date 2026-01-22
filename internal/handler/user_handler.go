package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"wavefy-be/helper"
	"wavefy-be/internal/dto"
	"wavefy-be/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Create godoc
// @Summary      Create user
// @Description  Create a new user with email and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateUserRequest true "Create user"
// @Success      200 {object} helper.Response{data=dto.UserResponse}
// @Failure      400 {object} helper.Response
// @Failure      409 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Create(c.Request.Context(), service.CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case err == service.ErrInvalidInput:
			helper.RespondError(c, http.StatusBadRequest, err.Error())
		case err == service.ErrEmailExists:
			helper.RespondError(c, http.StatusConflict, err.Error())
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondOK(c, dto.UserResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	})
}
