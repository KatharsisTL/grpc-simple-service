package externalserver

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
)

type Service interface {
	Hello(name string, withIdleSeconds uint64) (string, error)
}

type Server interface {
	Init(bindAddr string)
	SetService(svc Service)
	GetServer() *http.Server
}

type server struct {
	ctx    context.Context
	router *gin.Engine
	logger zerolog.Logger
	svc    Service
	http   *http.Server
}

func New(ctx context.Context, logger zerolog.Logger) *server {
	s := &server{
		router: gin.Default(),
		logger: logger,
		ctx:    ctx,
	}

	return s
}

func (s *server) SetService(svc Service) {
	s.svc = svc
	s.configureRouter()
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) Init(bindAddr string) {
	s.http = &http.Server{
		Addr:    bindAddr,
		Handler: s,
	}
}

func (s *server) GetServer() *http.Server {
	return s.http
}

func (s *server) configureRouter() {
	s.router.POST("/hello", func(ctx *gin.Context) {
		name := ctx.Query("name")
		idleString := ctx.Query("idle")

		var idle uint64 = 0
		var err error
		if idleString != "" {
			idle, err = strconv.ParseUint(idleString, 10, 64)
			if err != nil {
				_ = ctx.AbortWithError(http.StatusInternalServerError, err)

				return
			}
		}

		greeting, err := s.svc.Hello(name, idle)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, struct {
				Error string `json:"error"`
			}{
				Error: err.Error(),
			})

			return
		}

		ctx.String(http.StatusOK, greeting)
	})
}
