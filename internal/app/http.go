package app

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"wavefy-be/config"
	"wavefy-be/internal/handler"
)

// NewHTTP khởi tạo router.
func NewHTTP(db *gorm.DB, authCfg config.AuthConfig) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	h := handler.New(db)
	api := r.Group("/api")
	api.GET("/health", h.Health)
	api.GET("/db/ping", h.DBPing)
	registerUserRoutes(api, db)
	registerAuthRoutes(api, db, authCfg)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
