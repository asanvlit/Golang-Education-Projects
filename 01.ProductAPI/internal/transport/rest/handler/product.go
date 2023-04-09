package handler

import (
	"05.TestProductAPI/internal/core"
	"context"
	"fmt"
	"net/http"
	"time"
)
import "github.com/gofiber/fiber/v2"

type ProductService interface {
	GetAll(ctx context.Context) ([]*core.Product, error)
	GetById(ctx context.Context, id string) (*core.Product, error)
	CreateProduct(ctx context.Context, product *core.Product) (*core.Product, error)
}

type ProductHandler struct {
	service ProductService
}

func NewProductHandler(service ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (handler *ProductHandler) InitRoutes(app *fiber.App) {
	app.Get("/products", handler.GetAll)
	app.Get("/products/:productId", handler.GetById)
	app.Post("/products", handler.CreateProduct)
}

// GetAll
// @Summary Get all products
// @Tags products
// @Description Returns the list of the products
// @Produce json
// @Status 200
// @Router /products [get]
func (handler *ProductHandler) GetAll(ctx *fiber.Ctx) error {
	ctxTimeOut, cancel := context.WithTimeout(ctx.UserContext(), time.Second*5)
	defer cancel()

	productsChannel := make(chan []*core.Product, 0)

	var err error
	var products []*core.Product

	go func(channel chan<- []*core.Product) {
		products, err = handler.service.GetAll(ctxTimeOut)
		channel <- products
	}(productsChannel)

	if err != nil {
		return err
	}

	select {
	case <-ctxTimeOut.Done():
		fmt.Println("Processing timeout in Handler")
		break
	case products = <-productsChannel:
		fmt.Println("Finished processing in Handler")
	}

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			})
	}

	return ctx.Status(http.StatusOK).JSON(
		fiber.Map{
			"products": products,
		})
}

// GetById
// @Summary Get product by id
// @Tags products
// @Description Returns the product with this id
// @Produce json
// @Status 200
// @Param product-id path string true "Product ID"
// @Router /products/{product-id} [get]
func (handler *ProductHandler) GetById(ctx *fiber.Ctx) error {
	ctxTimeOut, cancel := context.WithTimeout(ctx.UserContext(), time.Second*2)
	defer cancel()

	productChannel := make(chan *core.Product, 0)

	var err error
	var product *core.Product

	go func(channel chan<- *core.Product) {
		product, err = handler.service.GetById(ctxTimeOut, ctx.Params("productId"))
		channel <- product
	}(productChannel)

	if err != nil {
		return err
	}

	select {
	case <-ctxTimeOut.Done():
		fmt.Println("Processing timeout in Handler")
		break
	case product = <-productChannel:
		fmt.Println("Finished processing in Handler")
	}

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			})
	}

	return ctx.Status(http.StatusOK).JSON(
		fiber.Map{
			"product": product,
		})
}

// CreateProduct
// @Summary Create new product
// @Tags products
// @Description Creates new product
// @Produce json
// @Status 201
// @Router /products [post]
func (handler *ProductHandler) CreateProduct(ctx *fiber.Ctx) error {
	product := &core.Product{}

	if err := ctx.BodyParser(product); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			})
	}

	ctxTimeOut, cancel := context.WithTimeout(ctx.UserContext(), time.Second*2)
	defer cancel()

	productChannel := make(chan *core.Product, 0)

	var err error
	var savedProduct *core.Product

	go func(channel chan<- *core.Product) {
		savedProduct, err = handler.service.CreateProduct(ctxTimeOut, product)
		channel <- savedProduct
	}(productChannel)

	if err != nil {
		return err
	}

	select {
	case <-ctxTimeOut.Done():
		fmt.Println("Processing timeout in Handler")
		break
	case savedProduct = <-productChannel:
		fmt.Println("Finished processing in Handler")
	}

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			})
	}

	return ctx.Status(http.StatusCreated).JSON(
		fiber.Map{
			"product": savedProduct,
		})
}
