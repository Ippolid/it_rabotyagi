package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(200, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string, err error) {
	response := Response{
		Success: false,
		Error:   message,
	}

	// В development режиме показываем детали ошибки
	if gin.Mode() == gin.DebugMode && err != nil {
		response.Error = err.Error()
	}

	c.JSON(statusCode, response)
}
