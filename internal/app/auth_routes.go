package app

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"wavefy-be/config"
	"wavefy-be/internal/handler"
	"wavefy-be/internal/repository"
	"wavefy-be/internal/service"
	"wavefy-be/internal/token"
)

func registerAuthRoutes(rg *gin.RouterGroup, db *gorm.DB, redisClient *redis.Client, cfg config.AuthConfig) {
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	userService := service.NewUserService(userRepo, roleRepo)
	refreshStore := token.NewRefreshTokenStore(redisClient, cfg.RefreshTokenSecret, cfg.RefreshTokenTTL)
	authService := service.NewAuthService(userService, userRepo, roleRepo, refreshStore, cfg)
	authHandler := handler.NewAuthHandler(authService, cfg)

	rg.POST("/auth/register", authHandler.Register)
	rg.POST("/auth/login", authHandler.Login)
	rg.POST("/auth/refresh", authHandler.Refresh)
	rg.POST("/auth/logout", authHandler.Logout)
}
