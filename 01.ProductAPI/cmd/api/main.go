package main

import (
	"05.TestProductAPI/internal/config"
	"05.TestProductAPI/internal/repostitory/mongo"
	"05.TestProductAPI/internal/service"
	"05.TestProductAPI/internal/transport/rest/handler"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"log"
	"time"
)

// @title 04.API-Product-Store-with-Mongo
// @version 0.2
// @description HW #4
// @termsOfService http://swagger.io/terms/

// @host localhost:3000
// @BasePath /
// @schemes http
func main() {
	if err := SetupViper(); err != nil {
		log.Fatal(err.Error())
	}

	app := fiber.New()

	config.SetupSwagger(app)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	mongoDataBase, err := config.SetupMongoDataBase(ctx, cancel)

	if err != nil {
		log.Fatal(err.Error())
	}

	productRepository := mongo.NewProductRepository(mongoDataBase.Collection("products"))
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	productHandler.InitRoutes(app)

	port := viper.GetString("http.port")
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

func SetupViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
