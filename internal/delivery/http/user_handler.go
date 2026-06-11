package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gokhan/orderly/internal/domain"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
}

func NewUserHandler(useCase domain.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: useCase,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req domain.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := domain.Validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.userUseCase.CreateUser(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) LoginUser(c *fiber.Ctx) error {
	var req domain.LoginUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := domain.Validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.userUseCase.LoginUser(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
