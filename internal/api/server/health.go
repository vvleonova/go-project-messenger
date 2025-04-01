package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// проверка подключения к БД
func (s Server) healthCheck(c *gin.Context) {
	if err := s.db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "database is not reachable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "database is connected"})
}
