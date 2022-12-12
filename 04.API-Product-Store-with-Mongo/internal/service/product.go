package service

import (
	"04.API-Product-Store-with-Mongo/internal/core"
	"context"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]*core.Product, error)
	GetById(ctx context.Context, id string) (*core.Product, error)
	Save(ctx context.Context, product *core.Product) (*core.Product, error)
}

type ProductService struct {
	productRepository ProductRepository
}

func NewProductService(repository ProductRepository) *ProductService {
	return &ProductService{productRepository: repository}
}

func (service *ProductService) GetAll(ctx context.Context) ([]*core.Product, error) {
	return service.productRepository.GetAll(ctx)
}

func (service *ProductService) GetById(ctx context.Context, id string) (*core.Product, error) {
	return service.productRepository.GetById(ctx, id)
}

func (service *ProductService) CreateProduct(ctx context.Context, product *core.Product) (*core.Product, error) {
	return service.productRepository.Save(ctx, product)
}
