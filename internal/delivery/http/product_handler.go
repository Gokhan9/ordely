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

// CreateProduct godoc
// @Summary Create a new product
// @Description Add a new product to the catalog (Admin only)
// @Tags products
// @Accept  json
// @Produce  json
// @Param product body domain.CreateProductRequest true "Product Info"
// @Success 201 {object} domain.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [post]
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

// GetProduct godoc
// @Summary Get product detail
// @Description Get detailed information about a product by ID
// @Tags products
// @Produce  json
// @Param id path int true "Product ID"
// @Success 200 {object} domain.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
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

// ListProducts godoc
// @Summary List products
// @Description Get a paginated list of products
// @Tags products
// @Produce  json
// @Param page_id query int true "Page ID" minimum(1)
// @Param page_size query int true "Page Size" minimum(5) maximum(20)
// @Success 200 {array} domain.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [get]
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
