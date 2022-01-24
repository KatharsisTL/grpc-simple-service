package internalserver

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"net/http"
)

type Server interface {
	Start(bindAddr string)
	Stop() error
	GetServer() *http.Server
}

type Service interface {
	IsReady() bool
}

type internalHTTP struct {
	router  *gin.Engine
	server  *http.Server
	service Service
	logger  zerolog.Logger
}

func New(logger zerolog.Logger, service Service) *internalHTTP {
	router := gin.New()
	router.Use(gin.Recovery())
	internal := &internalHTTP{
		router:  router,
		logger:  logger,
		service: service,
	}
	internal.registerHandlers()
	return internal
}

func (h *internalHTTP) registerHandlers() {
	h.router.Handle(http.MethodGet, "/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	h.router.Handle(http.MethodGet, "/readiness", func(ctx *gin.Context) {
		if h.service == nil {
			ctx.JSON(http.StatusInternalServerError, "not ready yet")
			return
		}

		if !h.service.IsReady() {
			ctx.JSON(http.StatusInternalServerError, "not ready yet")
			return
		}
		ctx.Status(http.StatusOK)
	})

	h.router.Handle(http.MethodGet, "/metrics", gin.WrapH(promhttp.Handler()))
}

func (h *internalHTTP) Start(bindURL string) {
	h.server = &http.Server{
		Addr:    bindURL,
		Handler: h.router,
	}
}

func (h *internalHTTP) GetServer() *http.Server {
	return h.server
}

func (h *internalHTTP) Stop() error {
	return nil
}
