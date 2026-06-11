package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gokhan/orderly/internal/domain"
	"github.com/gokhan/orderly/pkg/utils"
)

type OrderHandler struct {
	orderUseCase domain.OrderUseCase
	store        domain.UserUseCase // To get user ID from username
}

func NewOrderHandler(useCase domain.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: useCase,
	}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Place a new order with multiple products
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body domain.CreateOrderRequest true "Order Info"
// @Success 201 {object} domain.OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
// @Security BearerAuth
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	payload := c.Locals(authorizationPayloadKey).(*utils.Payload)
	
	var req domain.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := domain.Validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	order, err := h.orderUseCase.CreateOrder(c.Context(), payload.UserID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}
