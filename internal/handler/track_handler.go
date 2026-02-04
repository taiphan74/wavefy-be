package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"wavefy-be/helper"
	"wavefy-be/internal/dto"
	"wavefy-be/internal/model"
	"wavefy-be/internal/service"
)

type TrackHandler struct {
	service service.TrackService
}

func NewTrackHandler(service service.TrackService) *TrackHandler {
	return &TrackHandler{service: service}
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
