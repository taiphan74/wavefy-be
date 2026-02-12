package dto

type PresignTrackPutRequest struct {
	ContentType  string `json:"content_type"`
	ExpiresInSec *int   `json:"expires_in_sec"`
}

type PresignPutResponse struct {
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	ExpiresAt string            `json:"expires_at"`
	Key       string            `json:"key"`
	Bucket    string            `json:"bucket"`
}

type PresignGetRequest struct {
	Key          string `json:"key" binding:"required"`
	ExpiresInSec *int   `json:"expires_in_sec"`
}

type PresignGetResponse struct {
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	ExpiresAt string            `json:"expires_at"`
	Key       string            `json:"key"`
	Bucket    string            `json:"bucket"`
}

type DeleteObjectRequest struct {
	Key string `json:"key" binding:"required"`
}

type DeleteObjectResponse struct {
	Key     string `json:"key"`
	Bucket  string `json:"bucket"`
	Deleted bool   `json:"deleted"`
}
