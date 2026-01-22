package helper

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Time   string      `json:"time"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func RespondOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status: "ok",
		Code:   http.StatusOK,
		Time:   time.Now().UTC().Format(time.RFC3339),
		Data:   data,
	})
}

func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Status: "error",
		Code:   status,
		Time:   time.Now().UTC().Format(time.RFC3339),
		Error:  message,
	})
}
