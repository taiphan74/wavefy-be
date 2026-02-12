package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"wavefy-be/helper"
)

const (
	loginRateLimitMax = int64(5)
	loginRateLimitTTL = 60 * time.Second
	loginRateLimitKey = "login:ip:"
)

func LoginRateLimit(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if client == nil {
			c.Next()
			return
		}

		ip := strings.TrimSpace(c.ClientIP())
		if ip == "" {
			c.Next()
			return
		}

		key := loginRateLimitKey + ip
		count, err := client.Incr(c.Request.Context(), key).Result()
		if err != nil {
			helper.RespondError(c, http.StatusInternalServerError, "internal error")
			c.Abort()
			return
		}

		if count == 1 {
			if err := client.Expire(c.Request.Context(), key, loginRateLimitTTL).Err(); err != nil {
				helper.RespondError(c, http.StatusInternalServerError, "internal error")
				c.Abort()
				return
			}
		}

		if count >= loginRateLimitMax {
			helper.RespondError(c, http.StatusTooManyRequests, "too many login attempts")
			c.Abort()
			return
		}

		c.Next()
	}
}
