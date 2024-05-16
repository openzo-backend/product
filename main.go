package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/product/config"
	handlers "github.com/tanush-128/openzo_backend/product/internal/api"
	"github.com/tanush-128/openzo_backend/product/internal/middlewares"
	"github.com/tanush-128/openzo_backend/product/internal/pb"
	"github.com/tanush-128/openzo_backend/product/internal/repository"
	"github.com/tanush-128/openzo_backend/product/internal/service"
	"google.golang.org/grpc"
)

var UserClient pb.UserServiceClient

type User2 struct {
}

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config: %w", err))
	}

	db, err := connectToDB(cfg) // Implement database connection logic
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to database: %w", err))
	}

	// Initialize gRPC server
	// grpcServer := grpc.NewServer()

	// reflection.Register(grpcServer) // Optional for server reflection

	//Initialize gRPC client
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)
	UserClient = c

	imageConn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer imageConn.Close()
	imageClient := pb.NewImageServiceClient(imageConn)

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository, imageClient)

	go service.GrpcServer(cfg, &service.Server{
		ProductRepository: productRepository,
	})
	// Initialize HTTP server with Gin
	router := gin.Default()
	handler := handlers.NewHandler(&productService)

	router.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// router.Use(middlewares.JwtMiddleware(c))
	router.POST("/", handler.CreateProduct)
	router.GET("/store/:id", handler.GetProductsByStoreID)
	router.GET("/:id", handler.GetProductByID)
	router.Use(middlewares.NewMiddleware(c).JwtMiddleware)
	router.PUT("/:id", handler.UpdateProduct)
	router.DELETE("/:id", handler.DeleteProduct)

	// router.Use(middlewares.JwtMiddleware)

	router.Run(fmt.Sprintf(":%s", cfg.HTTPPort))

}
