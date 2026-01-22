package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"wavefy-be/helper"
)

// DBPing godoc
// @Summary      Ping database
// @Description  Check if database is reachable
// @Tags         db
// @Accept       json
// @Produce      json
// @Success      200 {object} helper.Response
// @Failure      503 {object} helper.Response
// @Router       /db/ping [get]
func (h *Handler) DBPing(c *gin.Context) {
	if h.db == nil {
		helper.RespondError(c, http.StatusServiceUnavailable, "db not configured")
		return
	}

	sqlDB, err := h.db.DB()
	if err != nil {
		helper.RespondError(c, http.StatusServiceUnavailable, err.Error())
		return
	}

	if err := sqlDB.PingContext(c.Request.Context()); err != nil {
		helper.RespondError(c, http.StatusServiceUnavailable, err.Error())
		return
	}

	helper.RespondOK(c, gin.H{"message": "ok"})
}
