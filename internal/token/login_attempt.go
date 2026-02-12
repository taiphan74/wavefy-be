package token

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultLoginAttemptTTL = 10 * time.Minute
	defaultLoginLockTTL    = 15 * time.Minute
	defaultLoginMaxAttempt = int64(10)
)

type LoginAttemptStore interface {
	IsLocked(ctx context.Context, email string) (bool, error)
	RecordFailure(ctx context.Context, email string) (count int64, locked bool, err error)
	Reset(ctx context.Context, email string) error
}

type loginAttemptStore struct {
	client        *redis.Client
	attemptTTL    time.Duration
	lockTTL       time.Duration
	maxAttempts   int64
	attemptPrefix string
	lockPrefix    string
}

func NewLoginAttemptStore(client *redis.Client, attemptTTL, lockTTL time.Duration, maxAttempts int64) LoginAttemptStore {
	if attemptTTL <= 0 {
		attemptTTL = defaultLoginAttemptTTL
	}
	if lockTTL <= 0 {
		lockTTL = defaultLoginLockTTL
	}
	if maxAttempts <= 0 {
		maxAttempts = defaultLoginMaxAttempt
	}

	return &loginAttemptStore{
		client:        client,
		attemptTTL:    attemptTTL,
		lockTTL:       lockTTL,
		maxAttempts:   maxAttempts,
		attemptPrefix: "login:mail:attempt:",
		lockPrefix:    "login:mail:lock:",
	}
}

func (s *loginAttemptStore) IsLocked(ctx context.Context, email string) (bool, error) {
	if s.client == nil {
		return false, nil
	}
	email = normalizeEmail(email)
	if email == "" {
		return false, errors.New("invalid email")
	}

	exists, err := s.client.Exists(ctx, s.lockKey(email)).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (s *loginAttemptStore) RecordFailure(ctx context.Context, email string) (int64, bool, error) {
	if s.client == nil {
		return 0, false, nil
	}
	email = normalizeEmail(email)
	if email == "" {
		return 0, false, errors.New("invalid email")
	}

	key := s.attemptKey(email)
	count, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, false, err
	}
	if count == 1 {
		if err := s.client.Expire(ctx, key, s.attemptTTL).Err(); err != nil {
			return 0, false, err
		}
	}

	if count >= s.maxAttempts {
		if err := s.client.Set(ctx, s.lockKey(email), "1", s.lockTTL).Err(); err != nil {
			return count, false, err
		}
		_ = s.client.Del(ctx, key).Err()
		return count, true, nil
	}

	return count, false, nil
}

func (s *loginAttemptStore) Reset(ctx context.Context, email string) error {
	if s.client == nil {
		return nil
	}
	email = normalizeEmail(email)
	if email == "" {
		return errors.New("invalid email")
	}
	return s.client.Del(ctx, s.attemptKey(email), s.lockKey(email)).Err()
}

func (s *loginAttemptStore) attemptKey(email string) string {
	return s.attemptPrefix + email
}

func (s *loginAttemptStore) lockKey(email string) string {
	return s.lockPrefix + email
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
