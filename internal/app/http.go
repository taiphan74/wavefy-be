package app

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"wavefy-be/internal/handler"
)

// RunHTTP khởi tạo router và chạy HTTP server.
func RunHTTP(addr string, db *sql.DB) error {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	h := handler.New(db)
	api := r.Group("/api")
	api.GET("/health", h.Health)
	api.GET("/db/ping", h.DBPing)

	return r.Run(addr)
}
