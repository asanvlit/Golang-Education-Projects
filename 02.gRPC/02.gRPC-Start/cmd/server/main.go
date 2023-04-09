package main

import (
	"02.gRPC-Start/internal/config"
	"02.gRPC-Start/internal/repository/mongo"
	"02.gRPC-Start/internal/server"
	pb "02.gRPC-Start/proto"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	host = "localhost"
	port = "5000"
)

func main() {
	ctx := context.Background()

	err := setupViper()

	if err != nil {
		log.Fatalf("error reading yml file: %v", err)
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("error starting tcp listener: %v", err)
	}

	mongoDataBase, err := config.SetupMongoDataBase(ctx)

	if err != nil {
		log.Fatalf("error starting mongo : %v", err)
	}

	productRepository := mongo.NewProductRepository(mongoDataBase.Collection("products"))

	productServer := server.NewProductServer(productRepository)

	grpcServer := grpc.NewServer()

	pb.RegisterProductServiceServer(grpcServer, productServer)

	log.Printf("gRPC started at %v\n", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("error starting gRPC : %v", err)
	}
}

func setupViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
