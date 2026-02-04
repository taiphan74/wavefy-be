package app

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"

	"wavefy-be/config"
	"wavefy-be/internal/handler"
	"wavefy-be/internal/service"
)

func registerUploadRoutes(rg *gin.RouterGroup, r2Client *s3.Client, r2Cfg config.R2Config) {
	uploadService := service.NewUploadService(r2Client, r2Cfg)
	uploadHandler := handler.NewUploadHandler(uploadService)

	rg.POST("/uploads/presign", uploadHandler.PresignPut)
}
