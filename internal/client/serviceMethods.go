package client

import (
	"context"
	"github.com/KatharsisTL/grpc-simple-service/pkg/grpcHello"
	"time"
)

func (s *service) Hello(name string, withIdleSeconds uint64) (string, error) {
	client, err := s.helloClient.Get()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	resp, err := client.Hello(ctx, &grpcHello.HelloRequest{
		Name:            name,
		WithIdleSeconds: withIdleSeconds,
	})
	defer cancel()
	if err != nil {
		return "", err
	}

	return resp.GetMessage(), nil
}
