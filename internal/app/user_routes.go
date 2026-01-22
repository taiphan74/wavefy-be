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

	rg.GET("/users", userHandler.List)
	rg.POST("/users", userHandler.Create)
	rg.GET("/users/:id", userHandler.Get)
	rg.PATCH("/users/:id", userHandler.Update)
	rg.DELETE("/users/:id", userHandler.Delete)
}
