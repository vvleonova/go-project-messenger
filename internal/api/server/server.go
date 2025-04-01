package server

import (
	"context"
	"fmt"
	"net/http"

	"go-messenger/internal/config"
	"go-messenger/internal/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// структура для управления HTTP-сервером
type Server struct {
	*http.Server
	db      storage
	service service
	logger  logger.Logger
}

// создание сервера
func New(config *config.Config, pg storage, service service, logger logger.Logger) *Server {
	// роутер
	r := gin.Default()

	// r := gin.New()
	// r.Use(gin.Recovery())
	// r.Use(middleware.RequestID()) - вспомогательная функция, перед тем как запрос в основной эндпоинт (есть ли пользователь в контактах перед отправкой сообщения)
	// r.Use(middleware.GinLogger(logger))
	// r.Use(middleware.GinContext())
	// r.Use(middleware.TotalRequestsCount())
	// r.Use(cors.Default())

	// инициализация структуры сервер
	server := &Server{
		&http.Server{
			Addr:    fmt.Sprintf(":%d", config.Port),
			Handler: r.Handler(), // все запросы попадают в этот обработчик
		},
		pg,
		service,
		logger,
	}

	// можно создать функцию init routers (в отдельной папке / в этой)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "welcome gin server")
	})

	// check health
	r.GET("/health/db", server.healthCheck)

	// user sign-up / sign-in
	r.POST("/user/sign-up", server.userSignUp)
	r.POST("/user/sign-in", server.userSignIn)

	// user authorization required
	authorized := r.Group("/user")
	authorized.Use(server.authRequired)
	{
		authorized.GET("/:phone", server.userGet)
		authorized.DELETE("/:phone", server.userDelete)
		authorized.PATCH("/:phone", server.userUpdate)
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
		s.logger.Error("exit reason", err)
	}
}
