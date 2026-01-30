package app

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"wavefy-be/config"
	"wavefy-be/internal/handler"
	"wavefy-be/internal/middleware"
)

// NewHTTP khởi tạo router.
func NewHTTP(db *gorm.DB, redisClient *redis.Client, authCfg config.AuthConfig) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	h := handler.New(db)
	api := r.Group("/api")
	api.GET("/health", h.Health)
	api.GET("/db/ping", h.DBPing)
	registerAuthRoutes(api, db, redisClient, authCfg)

	protected := api.Group("")
	protected.Use(middleware.JWTAuth(authCfg))
	registerUserRoutes(protected, db)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
