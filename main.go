package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/product/config"
	handlers "github.com/tanush-128/openzo_backend/product/internal/api"
	"github.com/tanush-128/openzo_backend/product/internal/pb"
	"github.com/tanush-128/openzo_backend/product/internal/repository"
	"github.com/tanush-128/openzo_backend/product/internal/service"
	"google.golang.org/grpc"
)

var UserClient pb.UserServiceClient

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config: %w", err))
	}

	db, err := connectToDB(cfg) // Implement database connection logic
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to database: %w", err))
	}

	conf := ReadConfig()
	p, _ := kafka.NewProducer(&conf)
	// topic := "notification"

	// go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	// Initialize gRPC server
	// grpcServer := grpc.NewServer()

	// reflection.Register(grpcServer) // Optional for server reflection

	// Initialize gRPC client
	// conn, err := grpc.Dial(cfg.UserGrpc, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer conn.Close()
	// c := pb.NewUserServiceClient(conn)
	// UserClient = c

	imageConn, err := grpc.Dial(cfg.ImageGrpc, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer imageConn.Close()
	imageClient := pb.NewImageServiceClient(imageConn)

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository, imageClient, p)

	go service.GrpcServer(cfg, &service.Server{
		ProductRepository: productRepository,
	})

	// Initialize Inventory Repository and Service
	inventoryTransactionRepository := repository.NewInventoryTransactionRepository(db)
	inventoryService := service.NewInventoryService(inventoryTransactionRepository)
	inventoryHandler := handlers.NewInventoryHandler(&inventoryService)

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
	router.GET("post/pincode/:pincode", handler.GetPostByPincode)
	router.GET("/:id", handler.GetProductByID)
	// router.Use(middlewares.NewMiddleware(c).JwtMiddleware)
	router.PUT("/:id", handler.ChangeProductQuantity)
	router.PUT("/", handler.UpdateProduct)
	router.PUT("/display_order/:id", handler.UpdateDisplayOrder)
	router.PUT("/display_order/batch", handler.BatchUpdateDisplayOrder)
	router.DELETE("/:id", handler.DeleteProduct)

	// Inventory routes
	router.POST("/inventory", inventoryHandler.CreateInventoryTransaction)
	router.GET("/inventory/:id", inventoryHandler.GetInventoryTransactionByID)
	router.PUT("/inventory/:id", inventoryHandler.UpdateInventoryTransaction)
	router.DELETE("/inventory/:id", inventoryHandler.DeleteInventoryTransaction)
	router.GET("/inventory/product/:product_id", inventoryHandler.GetAllTransactionsByProductID)

	// router.Use(middlewares.JwtMiddleware)

	router.Run(fmt.Sprintf(":%s", cfg.HTTPPort))

}
