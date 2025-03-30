package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/STLeee/test-grpc-microservice/common/protobuf"
)

type ServerA struct {
	protobuf.UnimplementedServiceAServer
	serviceBClient protobuf.ServiceBClient
}

type ServerAOption func(*ServerA) error

func WithServiceB(url string) ServerAOption {
	return func(s *ServerA) error {
		client, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("failed to create gRPC client: %v", err)
		}
		s.serviceBClient = protobuf.NewServiceBClient(client)
		return nil
	}
}

func NewServerA(opts ...ServerAOption) (*ServerA, error) {
	s := &ServerA{}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, fmt.Errorf("failed to apply option: %v", err)
		}
	}
	return s, nil
}

func (s *ServerA) GetData(ctx context.Context, req *protobuf.DataRequest) (*protobuf.DataResponse, error) {
	resB, err := s.serviceBClient.GetData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Service B: %v", err)
	}

	return &protobuf.DataResponse{Data: resB.Data}, nil
}

func main() {
	port := 50051 // Default port for Service A
	portStr := os.Getenv("SERVICE_A_PORT")
	if portStr != "" {
		port, _ = strconv.Atoi(portStr)
	}

	// Create service A with Service B client
	serviceBURL := os.Getenv("SERVICE_B_URL")
	if serviceBURL == "" {
		log.Fatal("SERVICE_B_URL environment variable is required")
	}
	serviceA, err := NewServerA(WithServiceB(serviceBURL))
	if err != nil {
		log.Fatalf("Failed to create Service A: %v", err)
	}

	// Create a listener on the specified port
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	protobuf.RegisterServiceAServer(grpcServer, serviceA)
	defer grpcServer.Stop()
	defer lis.Close()

	// Start the gRPC server
	log.Println("Service A is running on port", port)
	grpcServer.Serve(lis)
}
