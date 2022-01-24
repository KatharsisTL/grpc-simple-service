#!/bin/bash

# start evans to check grpc server
# evans pkg/grpcHello/hello.proto -p 8080

protoc -I ./pkg/grpc Hello --go_out=pkg/grpcHello --go_opt=paths=source_relative --go-grpc_out=pkg/grpcHello --go-grpc_opt=paths=source_relative pkg/grpcHello/hello.proto