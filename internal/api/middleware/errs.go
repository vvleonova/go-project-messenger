package middleware

import (
	"go-messenger/internal/api/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// выполнение хендлеров
		c.Next()

		// отсутствие ошибки
		err := c.Errors.Last()
		if err == nil {
			return
		}

		// обработка HTTPError
		if httpErr, ok := err.Err.(*errs.HTTPError); ok {
			c.JSON(httpErr.Status, gin.H{"error": httpErr.Message})
		} else {
			// неизвестная ошибка
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}
}
