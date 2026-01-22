package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"wavefy-be/internal/handler"
	"wavefy-be/internal/repository"
	"wavefy-be/internal/service"
)

func registerUserRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	rg.POST("/users", userHandler.Create)
}
