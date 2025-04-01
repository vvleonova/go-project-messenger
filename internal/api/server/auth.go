package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s Server) authRequired(c *gin.Context) {
	// получение токена из заголовка Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication token is required"})
		c.Abort()
		return
	}

	// Authorization: Bearer <token>
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication format"})
		c.Abort()
		return
	}

	tokenValue := tokenParts[1]

	// проверка токена
	claims, err := s.service.TokenVerify(tokenValue)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.Set("phone", claims["phone"])
	c.Next()
}
