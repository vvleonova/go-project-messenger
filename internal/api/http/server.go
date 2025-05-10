package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"go-messenger/internal/api/middleware"
	"go-messenger/internal/config"
	"go-messenger/internal/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// структура для управления HTTP-сервером
type Server struct {
	*http.Server
	service service
	logger  logger.Logger
	clients sync.Map
}

// создание сервера
func New(config *config.Config, service service, logger logger.Logger) *Server {
	// роутер
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandler())

	// инициализация структуры сервера
	server := &Server{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.Port),
			Handler: r.Handler(), // все запросы попадают в этот обработчик
		},
		service: service,
		logger:  logger,
	}

	// можно создать функцию init routers (в отдельной папке / в этой)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "welcome gin server")
	})

	// check health
	r.GET("/health", server.healthCheck)

	// user sign-up / sign-in
	authorization := r.Group("")
	{
		authorization.POST("/sign-up", server.userSignUp)
		authorization.POST("/sign-in", server.userSignIn)
		authorization.POST("/refresh", server.userRefreshToken)
	}

	// user authorization and identity required
	authorized := r.Group("/user")
	authorized.Use(middleware.AuthRequired(config.SecretKey))
	authorized.Use(middleware.IdentityRequired(server.service))
	{
		authorized.GET("/:phone", server.userGetPhone)
		authorized.DELETE("/:phone", server.userDelete)
		authorized.PATCH("/:phone", server.userUpdate)
		authorization.POST("/logout", server.userLogout)
		authorized.GET("/:phone/chats", server.userGetChats)
		authorized.GET("/:phone/chat", server.userGetChatMessages)
	}

	// user authorization websocket required
	websocket := r.Group("/user")
	websocket.Use(middleware.AuthWebsocketRequired(config.SecretKey))
	{
		websocket.GET("/ws", server.handleWebSocket)
	}

	return server
}

// запуск HTTP-сервера
func (s *Server) Start(ctx context.Context) {
	// загрузка логгера
	defer s.logger.Sync()

	// создание error group
	g, gCtx := errgroup.WithContext(ctx)

	// запуск HTTP-сервера в отдельной горутине
	g.Go(func() error {
		s.logger.Info("starting server", nil)
		return s.ListenAndServe()
	})

	// ожидание завершения извне в отдельной горутине
	g.Go(func() error {
		<-gCtx.Done()
		s.logger.Info("shutdown signal received", nil)
		return s.Shutdown(gCtx)
	})

	if err := g.Wait(); err != nil {
		s.logger.Info("exit reason", err)
	}
}
