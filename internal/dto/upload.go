package dto

type PresignPutRequest struct {
	Key          string `json:"key" binding:"required"`
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
