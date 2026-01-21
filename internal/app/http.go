package app

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"wavefy-be/internal/handler"
)

// NewHTTP khởi tạo router.
func NewHTTP(db *sql.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	h := handler.New(db)
	api := r.Group("/api")
	api.GET("/health", h.Health)
	api.GET("/db/ping", h.DBPing)

	return r
}
