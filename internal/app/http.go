package app

import (
	"time"

	"github.com/gin-contrib/cors"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
