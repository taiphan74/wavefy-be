package service

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"wavefy-be/config"
)

var ErrStorageNotConfigured = errors.New("storage not configured")

const (
	defaultPresignTTL = 15 * time.Minute
	minPresignTTL     = 1 * time.Minute
	maxPresignTTL     = 1 * time.Hour
)

type PresignPutInput struct {
	Key          string
	ContentType  string
	ExpiresInSec *int
}

type PresignPutOutput struct {
	URL       string
	Method    string
	Headers   map[string]string
	ExpiresAt time.Time
	Key       string
	Bucket    string
}

type UploadService interface {
	PresignPut(ctx context.Context, input PresignPutInput) (*PresignPutOutput, error)
}

type uploadService struct {
	presigner *s3.PresignClient
	bucket    string
}

func NewUploadService(r2Client *s3.Client, cfg config.R2Config) UploadService {
	var presigner *s3.PresignClient
	if r2Client != nil {
		presigner = s3.NewPresignClient(r2Client)
	}

	return &uploadService{
		presigner: presigner,
		bucket:    cfg.Bucket,
	}
}

func (s *uploadService) PresignPut(ctx context.Context, input PresignPutInput) (*PresignPutOutput, error) {
	if s.presigner == nil || strings.TrimSpace(s.bucket) == "" {
		return nil, ErrStorageNotConfigured
	}

	key := strings.TrimSpace(input.Key)
	if key == "" {
		return nil, ErrInvalidInput
	}

	contentType := strings.TrimSpace(input.ContentType)

	ttl := defaultPresignTTL
	if input.ExpiresInSec != nil {
		if *input.ExpiresInSec <= 0 {
			return nil, ErrInvalidInput
		}
		ttl = time.Duration(*input.ExpiresInSec) * time.Second
		if ttl < minPresignTTL || ttl > maxPresignTTL {
			return nil, ErrInvalidInput
		}
	}

	putInput := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	if contentType != "" {
		putInput.ContentType = aws.String(contentType)
	}

	presigned, err := s.presigner.PresignPutObject(ctx, putInput, func(o *s3.PresignOptions) {
		o.Expires = ttl
	})
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string, len(presigned.SignedHeader))
	for k, v := range presigned.SignedHeader {
		if len(v) > 0 {
			headers[http.CanonicalHeaderKey(k)] = v[0]
		}
	}

	return &PresignPutOutput{
		URL:       presigned.URL,
		Method:    presigned.Method,
		Headers:   headers,
		ExpiresAt: time.Now().UTC().Add(ttl),
		Key:       key,
		Bucket:    s.bucket,
	}, nil
}
