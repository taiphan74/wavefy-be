package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"wavefy-be/helper"
)

func (h *Handler) DBPing(c *gin.Context) {
	if h.db == nil {
		helper.RespondError(c, http.StatusServiceUnavailable, "db not configured")
		return
	}

	if err := h.db.PingContext(c.Request.Context()); err != nil {
		helper.RespondError(c, http.StatusServiceUnavailable, err.Error())
		return
	}

	helper.RespondOK(c, gin.H{"message": "ok"})
}
