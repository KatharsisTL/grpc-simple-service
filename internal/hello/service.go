package hello

import (
	"context"
	"github.com/KatharsisTL/grpc-simple-service/pkg/grpcHello"
	"github.com/KatharsisTL/grpc-simple-service/pkg/grpcStart"
	libHTTP "github.com/KatharsisTL/grpc-simple-service/pkg/http"
	"github.com/KatharsisTL/grpc-simple-service/pkg/internalserver"
	"github.com/KatharsisTL/grpc-simple-service/pkg/mutex"
	"github.com/KatharsisTL/grpc-simple-service/pkg/sig"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"os"
)

type Service interface {
	Start(ctx context.Context, g *errgroup.Group) error
	Stop() error
	IsReady() bool
}

type service struct {
	grpcHello.UnimplementedGrpcHelloServer

	ctx            context.Context
	cfg            *Config
	baseLogger     zerolog.Logger
	logger         zerolog.Logger
	internalServer internalserver.Server
	readyMutex     mutex.BoolMutex
	grpcServer     *grpc.Server
}

func New(ctx context.Context, g *errgroup.Group, cfg *Config) (Service, error) {
	logger := newLogger(cfg.LogLevel)

	svc := service{
		ctx:        ctx,
		cfg:        cfg,
		baseLogger: logger,
		logger:     logger.With().Str("component", "hello").Logger(),
		readyMutex: mutex.NewBoolMutex(false),
	}

	svc.internalServer = internalserver.New(
		logger.With().Str("component", "internal_http_server").Logger(),
		&svc)

	svc.internalServer.Start(cfg.ListenInternal)
	g.Go(libHTTP.MakeServerRunner(ctx, logger.With().Str("component", "internal_http_runner").Logger(), svc.internalServer.GetServer()))

	return &svc, nil
}

func (s *service) Start(ctx context.Context, g *errgroup.Group) error {
	g.Go(func() error {
		return sig.Listen(ctx)
	})

	// start grpc server
	s.grpcServer = grpc.NewServer()
	grpcHello.RegisterGrpcHelloServer(s.grpcServer, s)
	go func() {
		grpcStart.Start(s.grpcServer, s.cfg.Listen)
	}()
	g.Go(grpcStart.Start(s.grpcServer, s.cfg.Listen))

	s.setReady(true)
	return nil
}

func (s *service) Stop() error {
	// release resources
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

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
