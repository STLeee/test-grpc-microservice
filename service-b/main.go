package main

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"

	"github.com/STLeee/test-grpc-microservice/common/protobuf"
)

type ServerB struct {
	protobuf.UnimplementedServiceBServer
}

func NewServerB() (*ServerB, error) {
	return &ServerB{}, nil
}

func (s *ServerB) GetData(ctx context.Context, req *protobuf.DataRequest) (*protobuf.DataResponse, error) {
	data := "Hello, " + req.Key + " from B"
	return &protobuf.DataResponse{Data: data}, nil
}

func main() {
	port := 50052 // Default port for Service B
	portStr := os.Getenv("SERVICE_B_PORT")
	if portStr != "" {
		port, _ = strconv.Atoi(portStr)
	}

	// Create service B gRPC server
	serverB, err := NewServerB()
	if err != nil {
		log.Fatalf("Failed to create Service B: %v", err)
	}

	// Create a listener on the specified port
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	protobuf.RegisterServiceBServer(grpcServer, serverB)
	defer grpcServer.Stop()
	defer lis.Close()

	// Start the gRPC server
	log.Println("Service B is running on port", port)
	grpcServer.Serve(lis)
}
