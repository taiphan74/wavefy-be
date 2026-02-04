package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"wavefy-be/helper"
	"wavefy-be/internal/dto"
	"wavefy-be/internal/service"
)

type UploadHandler struct {
	service service.UploadService
}

func NewUploadHandler(service service.UploadService) *UploadHandler {
	return &UploadHandler{service: service}
}

// PresignPut godoc
// @Summary      Get presigned PUT URL
// @Tags         uploads
// @Accept       json
// @Produce      json
// @Param        request body dto.PresignPutRequest true "Presign PUT"
// @Success      200 {object} helper.Response{data=dto.PresignPutResponse}
// @Failure      400 {object} helper.Response
// @Failure      503 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /uploads/presign [post]
func (h *UploadHandler) PresignPut(c *gin.Context) {
	var req dto.PresignPutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	out, err := h.service.PresignPut(c.Request.Context(), service.PresignPutInput{
		Key:          req.Key,
		ContentType:  req.ContentType,
		ExpiresInSec: req.ExpiresInSec,
	})
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			helper.RespondError(c, http.StatusBadRequest, err.Error())
		case service.ErrStorageNotConfigured:
			helper.RespondError(c, http.StatusServiceUnavailable, err.Error())
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondOK(c, dto.PresignPutResponse{
		URL:       out.URL,
		Method:    out.Method,
		Headers:   out.Headers,
		ExpiresAt: out.ExpiresAt.Format(time.RFC3339),
		Key:       out.Key,
		Bucket:    out.Bucket,
	})
}
