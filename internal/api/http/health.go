package http

import (
	"go-messenger/internal/api/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// проверка подключения к БД
func (s *Server) healthCheck(c *gin.Context) {
	err := s.service.HealthCheck()
	if err != nil {
		c.Error(errs.ServiceUnavailable("database is unreachable"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "database is connected"})
}
