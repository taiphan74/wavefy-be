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
	globalRateLimitMax = int64(100)
	globalRateLimitTTL = time.Minute
	globalRateLimitKey = "rate:ip:"
)

func GlobalRateLimit(client *redis.Client) gin.HandlerFunc {
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

		key := globalRateLimitKey + ip
		count, err := client.Incr(c.Request.Context(), key).Result()
		if err != nil {
			helper.RespondError(c, http.StatusInternalServerError, "internal error")
			c.Abort()
			return
		}

		if count == 1 {
			if err := client.Expire(c.Request.Context(), key, globalRateLimitTTL).Err(); err != nil {
				helper.RespondError(c, http.StatusInternalServerError, "internal error")
				c.Abort()
				return
			}
		}

		if count > globalRateLimitMax {
			helper.RespondError(c, http.StatusTooManyRequests, "too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}
