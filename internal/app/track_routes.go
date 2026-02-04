package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"wavefy-be/internal/handler"
	"wavefy-be/internal/repository"
	"wavefy-be/internal/service"
)

func registerTrackRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	trackRepo := repository.NewTrackRepository(db)
	userRepo := repository.NewUserRepository(db)
	trackService := service.NewTrackService(trackRepo, userRepo)
	trackHandler := handler.NewTrackHandler(trackService)

	rg.GET("/tracks", trackHandler.List)
	rg.POST("/tracks", trackHandler.Create)
	rg.GET("/tracks/:id", trackHandler.Get)
	rg.PATCH("/tracks/:id", trackHandler.Update)
	rg.DELETE("/tracks/:id", trackHandler.Delete)
}
