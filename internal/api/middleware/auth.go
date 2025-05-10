package middleware

import (
	"go-messenger/internal/api/errs"
	"go-messenger/internal/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// получение токена из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(errs.Unauthorized("authentication token is required"))
			c.Abort()
			return
		}

		// Authorization: Bearer <token>
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Error(errs.Unauthorized("invalid authentication format"))
			c.Abort()
			return
		}

		tokenValue := tokenParts[1]

		// проверка токена
		claims, err := jwt.TokenVerify(tokenValue, secretKey)
		if err != nil {
			c.Error(errs.Unauthorized(err.Error()))
			c.Abort()
			return
		}

		c.Set("id", claims["id"])
		c.Next()
	}
}

func AuthWebsocketRequired(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// получение токена из запроса
		tokenValue := c.Query("token")
		if tokenValue == "" {
			c.Error(errs.Unauthorized("authentication token is required"))
			c.Abort()
			return
		}

		// проверка токена
		claims, err := jwt.TokenVerify(tokenValue, secretKey)
		if err != nil {
			c.Error(errs.Unauthorized(err.Error()))
			c.Abort()
			return
		}

		c.Set("id", claims["id"])
		c.Next()
	}
}
