package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"wavefy-be/helper"
	"wavefy-be/internal/dto"
	"wavefy-be/internal/model"
	"wavefy-be/internal/service"
)

type TrackHandler struct {
	service       service.TrackService
	uploadService service.UploadService
}

func NewTrackHandler(service service.TrackService, uploadService service.UploadService) *TrackHandler {
	return &TrackHandler{service: service, uploadService: uploadService}
}

const trackKeyPrefix = "tracks/"

func newTrackObjectKey(contentType string) string {
	base := uuid.NewString()
	ext := trackExtFromContentType(contentType)
	return trackKeyPrefix + base + ext
}

func trackExtFromContentType(contentType string) string {
	ct := strings.ToLower(strings.TrimSpace(contentType))
	switch ct {
	case "audio/mpeg", "audio/mp3":
		return ".mp3"
	case "audio/wav", "audio/x-wav":
		return ".wav"
	case "audio/flac":
		return ".flac"
	case "audio/aac":
		return ".aac"
	case "audio/ogg", "audio/opus":
		return ".ogg"
	case "audio/mp4", "audio/x-m4a", "audio/m4a":
		return ".m4a"
	case "audio/webm":
		return ".webm"
	default:
		return ""
	}
}

func normalizeTrackKey(key string) (string, error) {
	trimmed := strings.TrimSpace(key)
	if trimmed == "" || !strings.HasPrefix(trimmed, trackKeyPrefix) {
		return "", service.ErrInvalidInput
	}
	return trimmed, nil
}

// PresignPut godoc
// @Summary      Get presigned PUT URL for track audio
// @Tags         tracks
// @Accept       json
// @Produce      json
// @Param        request body dto.PresignTrackPutRequest true "Presign PUT"
// @Success      200 {object} helper.Response{data=dto.PresignPutResponse}
// @Failure      400 {object} helper.Response
// @Failure      503 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /tracks/audio/presign [post]
func (h *TrackHandler) PresignPut(c *gin.Context) {
	var req dto.PresignTrackPutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	out, err := h.uploadService.PresignPut(c.Request.Context(), service.PresignPutInput{
		Key:          newTrackObjectKey(req.ContentType),
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

// PresignGet godoc
// @Summary      Get presigned GET URL for track audio
// @Tags         tracks
// @Accept       json
// @Produce      json
// @Param        request body dto.PresignGetRequest true "Presign GET"
// @Success      200 {object} helper.Response{data=dto.PresignGetResponse}
// @Failure      400 {object} helper.Response
// @Failure      503 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /tracks/audio/presign-get [post]
func (h *TrackHandler) PresignGet(c *gin.Context) {
	var req dto.PresignGetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	key, err := normalizeTrackKey(req.Key)
	if err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	out, err := h.uploadService.PresignGet(c.Request.Context(), service.PresignGetInput{
		Key:          key,
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

	helper.RespondOK(c, dto.PresignGetResponse{
		URL:       out.URL,
		Method:    out.Method,
		Headers:   out.Headers,
		ExpiresAt: out.ExpiresAt.Format(time.RFC3339),
		Key:       out.Key,
		Bucket:    out.Bucket,
	})
}

// DeleteObject godoc
// @Summary      Delete uploaded track audio
// @Tags         tracks
// @Accept       json
// @Produce      json
// @Param        request body dto.DeleteObjectRequest true "Delete object"
// @Success      200 {object} helper.Response{data=dto.DeleteObjectResponse}
// @Failure      400 {object} helper.Response
// @Failure      503 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /tracks/audio/delete [post]
func (h *TrackHandler) DeleteObject(c *gin.Context) {
	var req dto.DeleteObjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	key, err := normalizeTrackKey(req.Key)
	if err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	out, err := h.uploadService.DeleteObject(c.Request.Context(), service.DeleteObjectInput{
		Key: key,
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

	helper.RespondOK(c, dto.DeleteObjectResponse{
		Key:     out.Key,
		Bucket:  out.Bucket,
		Deleted: true,
	})
}

// GetTrack godoc
// @Summary      Get track by id
// @Tags         tracks
// @Produce      json
// @Param        id path string true "Track ID"
// @Success      200 {object} helper.Response{data=dto.TrackResponse}
// @Failure      400 {object} helper.Response
// @Failure      404 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /tracks/{id} [get]
func (h *TrackHandler) Get(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	track, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrNotFound:
			helper.RespondError(c, http.StatusNotFound, err.Error())
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondOK(c, mapTrackResponse(track))
}

// ListTracks godoc
// @Summary      List tracks
// @Tags         tracks
// @Produce      json
// @Param        limit query int false "Limit" default(20)
// @Param        offset query int false "Offset" default(0)
// @Success      200 {object} helper.Response{data=[]dto.TrackResponse}
// @Failure      500 {object} helper.Response
// @Router       /tracks [get]
func (h *TrackHandler) List(c *gin.Context) {
	limit := parseIntQuery(c, "limit", 20)
	offset := parseIntQuery(c, "offset", 0)

	tracks, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.TrackResponse, 0, len(tracks))
	for i := range tracks {
		resp = append(resp, mapTrackResponse(&tracks[i]))
	}

	helper.RespondOK(c, resp)
}

// CreateTrack godoc
// @Summary      Create track
// @Tags         tracks
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateTrackRequest true "Create track"
// @Success      200 {object} helper.Response{data=dto.TrackResponse}
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /tracks [post]
func (h *TrackHandler) Create(c *gin.Context) {
	var req dto.CreateTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	track, err := h.service.Create(c.Request.Context(), service.CreateTrackInput{
		ArtistUserID: req.ArtistUserID,
		AlbumID:      req.AlbumID,
		Title:        req.Title,
		AudioURL:     req.AudioURL,
		DurationSec:  req.DurationSec,
		IsPublic:     req.IsPublic,
	})
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			helper.RespondError(c, http.StatusBadRequest, err.Error())
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondOK(c, mapTrackResponse(track))
}

// UpdateTrack godoc
// @Summary      Update track
// @Tags         tracks
// @Accept       json
// @Produce      json
// @Param        id path string true "Track ID"
// @Param        request body dto.UpdateTrackRequest true "Update track"
// @Success      200 {object} helper.Response{data=dto.TrackResponse}
// @Failure      400 {object} helper.Response
// @Failure      404 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /tracks/{id} [patch]
func (h *TrackHandler) Update(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	var req dto.UpdateTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	track, err := h.service.Update(c.Request.Context(), id, service.UpdateTrackInput{
		AlbumID:     req.AlbumID,
		Title:       req.Title,
		AudioURL:    req.AudioURL,
		DurationSec: req.DurationSec,
		IsPublic:    req.IsPublic,
	})
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			helper.RespondError(c, http.StatusBadRequest, err.Error())
		case service.ErrNotFound:
			helper.RespondError(c, http.StatusNotFound, err.Error())
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondOK(c, mapTrackResponse(track))
}

// DeleteTrack godoc
// @Summary      Delete track
// @Tags         tracks
// @Produce      json
// @Param        id path string true "Track ID"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      404 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /tracks/{id} [delete]
func (h *TrackHandler) Delete(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		switch err {
		case service.ErrNotFound:
			helper.RespondError(c, http.StatusNotFound, err.Error())
		default:
			helper.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondOK(c, gin.H{"deleted": true})
}

func mapTrackResponse(track *model.Track) dto.TrackResponse {
	var albumID *string
	if track.AlbumID != nil {
		value := track.AlbumID.String()
		albumID = &value
	}

	return dto.TrackResponse{
		ID:           track.ID.String(),
		ArtistUserID: track.ArtistUserID.String(),
		AlbumID:      albumID,
		Title:        track.Title,
		AudioURL:     track.AudioURL,
		DurationSec:  track.DurationSec,
		IsPublic:     track.IsPublic,
		PlayCount:    track.PlayCount,
		CreatedAt:    track.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    track.UpdatedAt.Format(time.RFC3339),
	}
}
