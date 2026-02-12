package app

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"wavefy-be/config"
	"wavefy-be/internal/handler"
	"wavefy-be/internal/repository"
	"wavefy-be/internal/service"
)

func registerTrackRoutes(rg *gin.RouterGroup, db *gorm.DB, r2Client *s3.Client, r2Cfg config.R2Config) {
	trackRepo := repository.NewTrackRepository(db)
	userRepo := repository.NewUserRepository(db)
	trackService := service.NewTrackService(trackRepo, userRepo)
	uploadService := service.NewUploadService(r2Client, r2Cfg)
	trackHandler := handler.NewTrackHandler(trackService, uploadService)

	rg.GET("/tracks", trackHandler.List)
	rg.POST("/tracks", trackHandler.Create)
	rg.GET("/tracks/:id", trackHandler.Get)
	rg.PATCH("/tracks/:id", trackHandler.Update)
	rg.DELETE("/tracks/:id", trackHandler.Delete)

	rg.POST("/tracks/audio/presign", trackHandler.PresignPut)
	rg.POST("/tracks/audio/presign-get", trackHandler.PresignGet)
	rg.POST("/tracks/audio/delete", trackHandler.DeleteObject)
}
