package http

import (
	"strings"

	"github.com/gin-gonic/gin"

	"go-messenger/internal/api/errs"
)

func (s Server) authRequired(c *gin.Context) {
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
	claims, err := s.service.TokenVerify(tokenValue)
	if err != nil {
		c.Error(errs.Unauthorized(err.Error()))
		c.Abort()
		return
	}

	c.Set("phone", claims["phone"])
	c.Next()
}
