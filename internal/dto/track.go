package dto

type CreateTrackRequest struct {
	ArtistUserID string  `json:"artist_user_id" binding:"required"`
	AlbumID      *string `json:"album_id"`
	Title        string  `json:"title" binding:"required"`
	AudioURL     string  `json:"audio_url" binding:"required"`
	ImageURL     *string `json:"image_url"`
	DurationSec  int     `json:"duration_sec" binding:"required"`
	IsPublic     *bool   `json:"is_public"`
}

type UpdateTrackRequest struct {
	AlbumID     *string `json:"album_id"`
	Title       *string `json:"title"`
	AudioURL    *string `json:"audio_url"`
	ImageURL    *string `json:"image_url"`
	DurationSec *int    `json:"duration_sec"`
	IsPublic    *bool   `json:"is_public"`
}

type TrackArtistResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type TrackResponse struct {
	ID          string              `json:"id"`
	Artist      TrackArtistResponse `json:"artist"`
	AlbumID     *string             `json:"album_id,omitempty"`
	Title       string              `json:"title"`
	AudioURL    string              `json:"audio_url"`
	ImageURL    *string             `json:"image_url,omitempty"`
	DurationSec int                 `json:"duration_sec"`
	IsPublic    bool                `json:"is_public"`
	PlayCount   int64               `json:"play_count"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
}
