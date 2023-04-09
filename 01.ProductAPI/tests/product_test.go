package tests

import (
	"05.TestProductAPI/internal/core"
	"05.TestProductAPI/internal/transport/rest/handler"
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]*core.Product, error)
	GetById(ctx context.Context, id string) (*core.Product, error)
	Save(ctx context.Context, product *core.Product) (*core.Product, error)
}

type MockProductRepository struct {
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{}
}

func (repository *MockProductRepository) GetAll(ctx context.Context) ([]*core.Product, error) {
	products := []*core.Product{
		{
			ID:          primitive.NewObjectID(),
			Name:        "Milk",
			Description: "Chocolate milk",
			Price:       88.2,
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "Bread",
			Description: "Grain bread",
			Price:       69,
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "Salt",
			Description: "Food salt",
			Price:       268.9,
		},
	}

	return products, nil
}

func (repository *MockProductRepository) GetById(ctx context.Context, id string) (*core.Product, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	product := &core.Product{
		ID:          objID,
		Name:        "Milk",
		Description: "Chocolate milk",
		Price:       88.2,
	}

	return product, nil
}

func (repository *MockProductRepository) Save(ctx context.Context, product *core.Product) (*core.Product, error) {
	product.ID = primitive.NewObjectID()
	return product, nil
}

type ProductService interface {
	GetAll(ctx context.Context) ([]*core.Product, error)
	GetById(ctx context.Context, id string) (*core.Product, error)
	CreateProduct(ctx context.Context, product *core.Product) (*core.Product, error)
}

type MockProductService struct {
	productRepository ProductRepository
}

func NewMockProductService(repository ProductRepository) *MockProductService {
	return &MockProductService{productRepository: repository}
}

func (service *MockProductService) GetAll(ctx context.Context) ([]*core.Product, error) {
	return service.productRepository.GetAll(ctx)
}

func (service *MockProductService) GetById(ctx context.Context, id string) (*core.Product, error) {
	return service.productRepository.GetById(ctx, id)
}

func (service *MockProductService) CreateProduct(ctx context.Context, product *core.Product) (*core.Product, error) {
	return service.productRepository.Save(ctx, product)
}

func TestGetAll(t *testing.T) {
	testCases := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "get HTTP status 200",
			route:        "/products",
			expectedCode: 200,
		},
		{
			description:  "get HTTP status 404",
			route:        "/invalid-path-products",
			expectedCode: 404,
		},
	}

	app := fiber.New()

	productRepository := NewMockProductRepository()
	productService := NewMockProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	productHandler.InitRoutes(app)

	for _, test := range testCases {
		req := httptest.NewRequest(http.MethodGet, test.route, nil)

		resp, _ := app.Test(req, 1)

		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}

func TestGetById(t *testing.T) {
	testCases := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "get HTTP status 200",
			route:        "/products",
			expectedCode: 200,
		},
	}

	app := fiber.New()

	productRepository := NewMockProductRepository()
	productService := NewMockProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	productHandler.InitRoutes(app)

	for _, test := range testCases {
		req := httptest.NewRequest(http.MethodGet, test.route+"/635976926e6524fe465bde94", nil)

		resp, _ := app.Test(req, 1)

		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}

func TestCreateProduct(t *testing.T) {
	testCases := []struct {
		description  string
		route        string
		product      *core.Product
		expectedCode int
	}{
		{
			description: "post HTTP status 201",
			route:       "/products",
			product: &core.Product{
				Name:        "Milk",
				Description: "Chocolate milk",
				Price:       88.2,
			},
			expectedCode: 201,
		},
		{
			description:  "post HTTP status 400",
			route:        "/products",
			expectedCode: 400,
		},
	}

	app := fiber.New()

	productRepository := NewMockProductRepository()
	productService := NewMockProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	productHandler.InitRoutes(app)

	for _, test := range testCases {
		var req *http.Request
		if test.product != nil {
			body, _ := json.Marshal(map[string]interface{}{
				"name":        test.product.Name,
				"description": test.product.Description,
				"price":       test.product.Price,
			})
			req = httptest.NewRequest(http.MethodPost, test.route, bytes.NewReader(body))
		} else {
			req = httptest.NewRequest(http.MethodPost, test.route, nil)
		}

		req.Header.Add("Content-Type", "application/json")

		resp, _ := app.Test(req, 1)

		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
