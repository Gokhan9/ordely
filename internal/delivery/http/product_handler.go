package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gokhan/orderly/internal/domain"
)

type ProductHandler struct {
	productUseCase domain.ProductUseCase
}

func NewProductHandler(useCase domain.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: useCase,
	}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req domain.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := domain.Validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	product, err := h.productUseCase.CreateProduct(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid product id"})
	}

	product, err := h.productUseCase.GetProduct(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	var req domain.ListProductsRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	products, err := h.productUseCase.ListProducts(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}
