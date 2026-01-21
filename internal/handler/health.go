package handler

import (
	"github.com/gin-gonic/gin"

	"wavefy-be/helper"
)

func (h *Handler) Health(c *gin.Context) {
	helper.RespondOK(c, gin.H{"message": "ok"})
}
