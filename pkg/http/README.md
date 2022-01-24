# HTTP package

## MakeServerRunner

See example:

```go
package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	libHTTP "git.wildberries.ru/courier/libraries/http"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
		log.Info().Str("signal", (<-shutdown).String()).Msg("signal received")
		cancel()
	}()

	r := http.NewServeMux()
	r.HandleFunc("/status", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(writer, "ok")
	})

	server := &http.Server{
		Addr:    ":8888",
		Handler: r}

	group, groupCtx := errgroup.WithContext(ctx)
    // groupCtx will expire when ctx expires or one of the jobs returns an error
	group.Go(libHTTP.MakeServerRunner(
		groupCtx,
		log.Logger.With().Str("http_server", "internal API").Logger(),
		server))
	group.Go(func() error {
		// another job
		return nil
	})

	if err := group.Wait(); err != nil {
		log.Error().Err(err).Msg("one of job has been failed")
	}
}

```