package service

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/tanush-128/openzo_backend/product/config"
	"github.com/tanush-128/openzo_backend/product/internal/pb"
	"github.com/tanush-128/openzo_backend/product/internal/repository"
	"google.golang.org/grpc"
)

type Server struct {
	pb.ProductServiceServer
	ProductRepository repository.ProductRepository
}

func GrpcServer(
	cfg *config.Config,
	server *Server,
) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Server listening at %v", lis.Addr())
	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, server)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

func (s *Server) ChangeProductQuantity(ctx context.Context, req *pb.ChangeProductQuantityRequest) (*pb.ChangeProductQuantityResponse, error) {
	// Implement your business logic here
	err := s.ProductRepository.ChangeProductQuantity(req.GetProductId(), int(req.GetQuantity()))
	if err != nil {
		return nil, err
	}
	return &pb.ChangeProductQuantityResponse{
		Status: "success",
	}, nil

}
