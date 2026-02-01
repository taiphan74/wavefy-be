package token

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type VerifyEmailTokenStore interface {
	Create(ctx context.Context, userID string) (string, error)
	Verify(ctx context.Context, token string) (string, error)
	Revoke(ctx context.Context, token string) error
}

type verifyEmailStore struct {
	client *redis.Client
	secret []byte
	ttl    time.Duration
	prefix string
}

func NewVerifyEmailTokenStore(client *redis.Client, secret string, ttl time.Duration) VerifyEmailTokenStore {
	return &verifyEmailStore{
		client: client,
		secret: []byte(secret),
		ttl:    ttl,
		prefix: "verify:",
	}
}

func (s *verifyEmailStore) Create(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", errors.New("invalid user id")
	}
	token := uuid.NewString()
	hash := s.sign(token, userID)
	value := userID + ":" + hash

	if err := s.client.Set(ctx, s.key(token), value, s.ttl).Err(); err != nil {
		return "", err
	}
	return token, nil
}

func (s *verifyEmailStore) Verify(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", errors.New("invalid token")
	}
	value, err := s.client.Get(ctx, s.key(token)).Result()
	if err != nil {
		return "", err
	}
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid token")
	}
	userID := parts[0]
	hash := parts[1]
	if !hmac.Equal([]byte(hash), []byte(s.sign(token, userID))) {
		return "", errors.New("invalid token")
	}
	return userID, nil
}

func (s *verifyEmailStore) Revoke(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	return s.client.Del(ctx, s.key(token)).Err()
}

func (s *verifyEmailStore) key(token string) string {
	return s.prefix + token
}

func (s *verifyEmailStore) sign(token, userID string) string {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(token))
	mac.Write([]byte(":"))
	mac.Write([]byte(userID))
	return hex.EncodeToString(mac.Sum(nil))
}
