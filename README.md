# grpc-simple-service
Simple gRPC client and server

Services with graceful shutdown, grpc client is lazy: connects to gRPC server in first client call

Start gRPC server hello

`go run cmd/hello/main.go`

Start gRPC client

`go run cmd/client/main.go`

Run http request

`curl -X POST "http://localhost:8081/hello?name=My_name&idle=5"`