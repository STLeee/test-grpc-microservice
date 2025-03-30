package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/STLeee/test-grpc-microservice/common/protobuf"
	"github.com/STLeee/test-grpc-microservice/service-api/docs"
)

func main() {
	host := os.Getenv("SERVICE_API_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := 8080
	portStr := os.Getenv("SERVICE_API_PORT")
	if portStr != "" {
		port, _ = strconv.Atoi(portStr)
	}

	serviceAUrl := os.Getenv("SERVICE_A_URL")
	if serviceAUrl == "" {
		serviceAUrl = fmt.Sprintf("%s:%d", host, 50051)
	}
	serviceBUrl := os.Getenv("SERVICE_B_URL")
	if serviceBUrl == "" {
		serviceBUrl = fmt.Sprintf("%s:%d", host, 50052)
	}

	// Create controller with gRPC clients
	controller, err := NewController(
		WithServiceA(serviceAUrl),
		WithServiceB(serviceBUrl),
	)
	if err != nil {
		log.Fatal("Failed to create controller:", err)
	}

	// Setup server
	engine := gin.Default()
	apiRouter := engine.Group("/api")
	apiRouter.GET("/service-a/data", controller.GetDataFromServiceA)
	apiRouter.GET("/service-b/data", controller.GetDataFromServiceB)

	// Swagger
	docs.SwaggerInfo.Title = "API Service"
	docs.SwaggerInfo.Description = "API Service for gRPC Microservices"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", host, port)
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run server
	engine.Run(fmt.Sprintf(":%d", port))
}

type Controller struct {
	serviceAClient protobuf.ServiceAClient
	serviceBClient protobuf.ServiceBClient
}

type ControllerOption func(*Controller) error

func WithServiceA(url string) ControllerOption {
	return func(c *Controller) error {
		client, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("failed to create gRPC client: %v", err)
		}
		c.serviceAClient = protobuf.NewServiceAClient(client)
		return nil
	}
}

func WithServiceB(url string) ControllerOption {
	return func(c *Controller) error {
		client, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("failed to create gRPC client: %v", err)
		}
		c.serviceBClient = protobuf.NewServiceBClient(client)
		return nil
	}
}

func NewController(opts ...ControllerOption) (*Controller, error) {
	c := &Controller{}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("failed to apply option: %v", err)
		}
	}
	return c, nil
}

// @Summary Get data
// @Description Get data from Service A
// @Tags service-a
// @Param key query string true "Key"
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /service-a/data [get]
func (c *Controller) GetDataFromServiceA(ctx *gin.Context) {
	key := ctx.Query("key")

	req := &protobuf.DataRequest{
		Key: key,
	}
	res, err := c.serviceAClient.GetData(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to call Service A"})
		return
	}
	ctx.JSON(200, res)
}

// @Summary Get data
// @Description Get data from Service B
// @Tags service-b
// @Param key query string true "Key"
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /service-b/data [get]
func (c *Controller) GetDataFromServiceB(ctx *gin.Context) {
	key := ctx.Query("key")

	req := &protobuf.DataRequest{
		Key: key,
	}
	res, err := c.serviceBClient.GetData(ctx, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to call Service B"})
		return
	}
	ctx.JSON(200, res)
}
