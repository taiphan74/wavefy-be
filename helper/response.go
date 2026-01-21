package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func RespondOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status: "ok",
		Data:   data,
	})
}

func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Status: "error",
		Error:  message,
	})
}
