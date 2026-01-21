package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DBPing(c *gin.Context) {
	if h.db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "error",
			"error":  "db not configured",
		})
		return
	}

	if err := h.db.PingContext(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
