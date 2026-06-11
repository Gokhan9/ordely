package http

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gokhan/orderly/pkg/utils"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorizationHeader := c.Get(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "authorization header is not provided"})
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header format"})
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unsupported authorization type"})
		}

		accessToken := fields[1]
		payload, err := utils.VerifyToken(accessToken, secretKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		c.Locals(authorizationPayloadKey, payload)
		return c.Next()
	}
}
