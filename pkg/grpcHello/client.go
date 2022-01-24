package grpcHello

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

type HelloClient struct {
	address        string
	client         GrpcHelloClient
	connectTimeout time.Duration
}

func New(address string, connectTimeout time.Duration) *HelloClient {
	return &HelloClient{address: address, connectTimeout: connectTimeout}
}

func (c *HelloClient) Get() (GrpcHelloClient, error) {
	if c.client != nil {
		return c.client, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.connectTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, c.address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("grpc client dial failed: %w", err)
	}

	c.client = NewGrpcHelloClient(conn)

	return c.client, nil
}
