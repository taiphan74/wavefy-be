package handler

import (
	"github.com/gin-gonic/gin"

	"wavefy-be/helper"
)

// Health godoc
// @Summary      Ping server
// @Description  Check if service is up
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200 {object} helper.Response
// @Router       /health [get]
func (h *Handler) Health(c *gin.Context) {
	helper.RespondOK(c, gin.H{"message": "ok"})
}
