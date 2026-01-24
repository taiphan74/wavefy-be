package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"wavefy-be/config"
	"wavefy-be/internal/handler"
	"wavefy-be/internal/repository"
	"wavefy-be/internal/service"
)

func registerAuthRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg config.AuthConfig) {
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	userService := service.NewUserService(userRepo, roleRepo)
	authService := service.NewAuthService(userService, userRepo, cfg)
	authHandler := handler.NewAuthHandler(authService)

	rg.POST("/auth/register", authHandler.Register)
	rg.POST("/auth/login", authHandler.Login)
}
