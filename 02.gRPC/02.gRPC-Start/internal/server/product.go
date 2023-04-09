package server

import (
	"02.gRPC-Start/internal/core"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
import pb "02.gRPC-Start/proto"

type ProductRepository interface {
	GetById(ctx context.Context, id string) (*core.Product, error)
}

type ProductServer struct {
	pb.ProductServiceServer
	productRepository ProductRepository
}

func NewProductServer(repository ProductRepository) *ProductServer {
	return &ProductServer{
		productRepository: repository,
	}
}

func (server *ProductServer) GetProduct(ctx context.Context, request *pb.ProductRequest) (response *pb.ProductResponse, err error) {
	product, err := server.productRepository.GetById(ctx, request.GetId())

	if err != nil {
		return nil, err
	}

	if product == nil {
		product = &core.Product{
			ID:          primitive.ObjectID{},
			Name:        "",
			Description: "",
			Price:       0,
		}
	}

	return &pb.ProductResponse{Name: product.Name}, nil
}
