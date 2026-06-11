package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Orderly API v1.0",
	})

	// Middlewares
	app.Use(logger.New())
	app.Use(recover.New())

	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Orderly API is running smoothly",
		})
	})

	log.Fatal(app.Listen(":8080"))
}
