package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"wavefy-be/config"
)

type AccessTokenClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func IssueAccessToken(cfg config.AuthConfig, subject, role string) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(cfg.AccessTokenTTL)
	claims := AccessTokenClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    cfg.AccessTokenIss,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}
