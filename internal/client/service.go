package client

import (
	"context"
	"github.com/KatharsisTL/grpc-simple-service/internal/client/externalserver"
	"github.com/KatharsisTL/grpc-simple-service/pkg/grpcHello"
	libHTTP "github.com/KatharsisTL/grpc-simple-service/pkg/http"
	"github.com/KatharsisTL/grpc-simple-service/pkg/internalserver"
	"github.com/KatharsisTL/grpc-simple-service/pkg/mutex"
	"github.com/KatharsisTL/grpc-simple-service/pkg/sig"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"os"
)

type Service interface {
	Start(ctx context.Context, g *errgroup.Group) error
	Stop() error
	IsReady() bool
}

type service struct {
	ctx            context.Context
	cfg            *Config
	baseLogger     zerolog.Logger
	logger         zerolog.Logger
	internalServer internalserver.Server
	externalServer externalserver.Server
	readyMutex     mutex.BoolMutex
	helloClient    *grpcHello.HelloClient
}

func New(ctx context.Context, g *errgroup.Group, cfg *Config) (Service, error) {
	logger := newLogger(cfg.LogLevel)

	svc := service{
		ctx:            ctx,
		cfg:            cfg,
		baseLogger:     logger,
		logger:         logger.With().Str("component", "client").Logger(),
		externalServer: externalserver.New(ctx, logger.With().Str("component", "external_http_server").Logger()),
		readyMutex:     mutex.NewBoolMutex(false),
		helloClient:    grpcHello.New(cfg.HelloClientUrl, cfg.HelloClientConnectTimeout),
	}

	svc.internalServer = internalserver.New(
		logger.With().Str("component", "internal_http_server").Logger(),
		&svc)

	svc.internalServer.Start(cfg.ListenInternal)
	g.Go(libHTTP.MakeServerRunner(ctx, logger.With().Str("component", "internal_http_runner").Logger(), svc.internalServer.GetServer()))

	svc.externalServer.SetService(&svc)

	return &svc, nil
}

func (s *service) Start(ctx context.Context, g *errgroup.Group) error {
	g.Go(func() error {
		return sig.Listen(ctx)
	})

	s.externalServer.Init(s.cfg.Listen)
	g.Go(libHTTP.MakeServerRunner(ctx, s.baseLogger.With().Str("component", "external_http_runner").Logger(), s.externalServer.GetServer()))

	s.setReady(true)
	return nil
}

func (s *service) Stop() error {
	// release resources

	return nil
}

func (s *service) setReady(isReady bool) {
	s.readyMutex.SetBoolValue(isReady)
}

func (s *service) IsReady() bool {
	return s.readyMutex.BoolValue()
}

func newLogger(logLevel string) zerolog.Logger {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.DebugLevel
	}

	return zerolog.New(os.Stdout).Level(level).With().Timestamp().Logger()
}
