package grpcStart

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Start(grpcServer *grpc.Server, listen string) func() error {
	return func() error {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

		l, err := net.Listen("tcp", listen)
		if err != nil {
			return fmt.Errorf("grpc listen failed: %w", err)
		}

		errChan := make(chan error)
		go func() {
			errChan <- grpcServer.Serve(l)
		}()

		select {
		case err := <-errChan:
			return err
		case <-sigint:
			return errors.New("sigterm")
		}
	}
}
