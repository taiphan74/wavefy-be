package app

import "github.com/gin-gonic/gin"

// RunHTTP khởi tạo router và chạy HTTP server.
func RunHTTP(addr string) error {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	api := r.Group("/api")
	_ = api

	return r.Run(addr)
}
