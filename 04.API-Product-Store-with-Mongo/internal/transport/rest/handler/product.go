package handler

import (
	"04.API-Product-Store-with-Mongo/internal/core"
	"context"
	"net/http"
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

func (handler *ProductHandler) GetAll(ctx *fiber.Ctx) error {
	products, err := handler.service.GetAll(ctx.UserContext())

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

func (handler *ProductHandler) GetById(ctx *fiber.Ctx) error {
	product, err := handler.service.GetById(ctx.UserContext(), ctx.Params("productId"))

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

func (handler *ProductHandler) CreateProduct(ctx *fiber.Ctx) error {
	product := &core.Product{}

	if err := ctx.BodyParser(product); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			})
	}

	savedProduct, err := handler.service.CreateProduct(ctx.UserContext(), product)

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
