package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"wavefy-be/config"
)

func JWTAuth(cfg config.AuthConfig) gin.HandlerFunc {
	secret := []byte(cfg.JWTSecret)

	return func(c *gin.Context) {
		tokenStr, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"code":   http.StatusUnauthorized,
				"error":  err.Error(),
			})
			return
		}

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"code":   http.StatusUnauthorized,
				"error":  "invalid token",
			})
			return
		}

		if cfg.AccessTokenIss != "" && claims.Issuer != cfg.AccessTokenIss {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"code":   http.StatusUnauthorized,
				"error":  "invalid token issuer",
			})
			return
		}

		c.Set("auth_subject", claims.Subject)
		c.Next()
	}
}

func extractBearerToken(value string) (string, error) {
	if value == "" {
		return "", errors.New("missing authorization header")
	}
	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid authorization header")
	}
	if parts[1] == "" {
		return "", errors.New("invalid authorization header")
	}
	return parts[1], nil
}
