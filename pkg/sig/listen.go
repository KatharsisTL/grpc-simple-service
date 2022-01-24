package sig

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

var ErrShutdownSignalReceived = errors.New("shutdown")

func Listen(ctx context.Context) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-shutdown:
		log.Info().Str("signal", s.String()).Msg("signal received")
		return ErrShutdownSignalReceived

	case <-ctx.Done():
		return nil
	}
}
