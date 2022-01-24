// Copyright © 2020 Wildberries. All rights reserved.

package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// MakeServerRunner creates closure for running HTTP server in separate goroutine
//
//goland:noinspection GoUnusedExportedFunction
func MakeServerRunner(ctx context.Context, logger zerolog.Logger, server *http.Server) func() error {
	return func() error {
		errCh := make(chan error)
		go func() {
			errCh <- server.ListenAndServe()
		}()
		logger.Info().Msg("http server has been started")
		select {
		case err := <-errCh:
			logger.Error().Err(err).Msg("http server failed")
			return err
		case <-ctx.Done():
			logger.Info().Msg("http server is stopping")
			shutdownCtxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
			err := server.Shutdown(shutdownCtxTimeout)
			cancel()
			if err == nil || err == http.ErrServerClosed {
				logger.Info().Msg("http server has been stopped")
				return nil
			}
			logger.Error().Err(err).Msg("failed to stop http server")
			return fmt.Errorf("shutdowning http server: %w", err)
		}
	}
}
