package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"wavefy-be/helper"
	"wavefy-be/internal/dto"
	"wavefy-be/internal/model"
	"wavefy-be/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Get godoc
// @Summary      Get user by id
// @Tags         users
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200 {object} helper.Response{data=dto.UserResponse}
// @Failure      400 {object} helper.Response
// @Failure      404 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrNotFound:
			helper.RespondError(c, http.StatusNotFound, err.Error())
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondOK(c, mapUserResponse(user))
}

// List godoc
// @Summary      List users
// @Tags         users
// @Produce      json
// @Param        limit query int false "Limit" default(20)
// @Param        offset query int false "Offset" default(0)
// @Success      200 {object} helper.Response{data=[]dto.UserResponse}
// @Failure      500 {object} helper.Response
// @Router       /users [get]
func (h *UserHandler) List(c *gin.Context) {
	limit := parseIntQuery(c, "limit", 20)
	offset := parseIntQuery(c, "offset", 0)

	users, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		resp = append(resp, mapUserResponse(&users[i]))
	}

	helper.RespondOK(c, resp)
}

// Update godoc
// @Summary      Update user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        request body dto.UpdateUserRequest true "Update user"
// @Success      200 {object} helper.Response{data=dto.UserResponse}
// @Failure      400 {object} helper.Response
// @Failure      404 {object} helper.Response
// @Failure      409 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /users/{id} [patch]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Update(c.Request.Context(), id, service.UpdateUserInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	})
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			helper.RespondError(c, http.StatusBadRequest, err.Error())
		case service.ErrEmailExists:
			helper.RespondError(c, http.StatusConflict, err.Error())
		case service.ErrNotFound:
			helper.RespondError(c, http.StatusNotFound, err.Error())
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondOK(c, mapUserResponse(user))
}

// Delete godoc
// @Summary      Delete user
// @Tags         users
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      404 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		switch err {
		case service.ErrNotFound:
			helper.RespondError(c, http.StatusNotFound, err.Error())
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondOK(c, gin.H{"deleted": true})
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

func mapUserResponse(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

func parseUUIDParam(c *gin.Context, key string) (uuid.UUID, error) {
	value := c.Param(key)
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, errors.New("invalid id")
	}
	return id, nil
}

func parseIntQuery(c *gin.Context, key string, fallback int) int {
	value := c.Query(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return fallback
	}
	return parsed
}
