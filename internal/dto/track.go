package dto

type CreateTrackRequest struct {
	ArtistUserID string  `json:"artist_user_id" binding:"required"`
	AlbumID      *string `json:"album_id"`
	Title        string  `json:"title" binding:"required"`
	AudioURL     string  `json:"audio_url" binding:"required"`
	DurationSec  int     `json:"duration_sec" binding:"required"`
	IsPublic     *bool   `json:"is_public"`
}

type UpdateTrackRequest struct {
	AlbumID     *string `json:"album_id"`
	Title       *string `json:"title"`
	AudioURL    *string `json:"audio_url"`
	DurationSec *int    `json:"duration_sec"`
	IsPublic    *bool   `json:"is_public"`
}

type TrackResponse struct {
	ID           string  `json:"id"`
	ArtistUserID string  `json:"artist_user_id"`
	AlbumID      *string `json:"album_id,omitempty"`
	Title        string  `json:"title"`
	AudioURL     string  `json:"audio_url"`
	DurationSec  int     `json:"duration_sec"`
	IsPublic     bool    `json:"is_public"`
	PlayCount    int64   `json:"play_count"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}
