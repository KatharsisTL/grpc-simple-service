package hello

import (
	"context"
	"errors"
	"github.com/KatharsisTL/grpc-simple-service/pkg/grpcHello"
	"time"
)

func (s *service) Hello(_ context.Context, helloRequest *grpcHello.HelloRequest) (*grpcHello.HelloResponse, error) {
	if helloRequest == nil {
		return &grpcHello.HelloResponse{}, errors.New("hello request is empty")
	}

	if helloRequest.GetWithIdleSeconds() > 0 {
		time.Sleep(time.Duration(helloRequest.GetWithIdleSeconds()) * time.Second)
	}

	name := "mr. Unnamed"
	if helloRequest != nil && helloRequest.GetName() != "" {
		name = helloRequest.GetName()
	}

	return &grpcHello.HelloResponse{
		Message: "Hello, " + name + "!",
	}, nil
}
