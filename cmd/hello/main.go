package main

import (
	"context"
	"errors"
	"github.com/KatharsisTL/grpc-simple-service/internal/hello"
	"github.com/KatharsisTL/grpc-simple-service/pkg/sig"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg, err := hello.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed get config")
	}

	g, ctx := errgroup.WithContext(context.Background())

	svc, err := hello.New(ctx, g, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("error create service")
	}

	err = svc.Start(ctx, g)
	if err != nil {
		log.Fatal().Err(err).Msg("error start service")
	}

	if err := g.Wait(); err != nil {
		if !errors.Is(err, sig.ErrShutdownSignalReceived) {
			log.Error().Err(err).Msg("errgroup error")
		}

		log.Info().Msg("service stopping")

		err = svc.Stop()
		if err != nil {
			log.Fatal().Err(err).Msg("error stop service")
		}

		log.Info().Msg("service stopped")
	}
}
