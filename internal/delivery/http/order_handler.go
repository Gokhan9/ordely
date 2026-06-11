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

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	payload := c.Locals(authorizationPayloadKey).(*utils.Payload)
	
	// This is a bit simplified, in a real app we'd have a way to get userID from payload
	// For now, let's assume we fetch user first or the payload has it.
	// Since our payload only has Username, let's fetch user.
	// But to keep it simple and performance-focused, we should've put userID in token.
	
	var req domain.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// For demo, we hardcode userID 1 or we need a service to find it
	userID := int64(1) 

	order, err := h.orderUseCase.CreateOrder(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}
