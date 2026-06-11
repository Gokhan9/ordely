package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gokhan/orderly/internal/delivery/http"
	"github.com/gokhan/orderly/internal/repository/db"
	"github.com/gokhan/orderly/internal/usecase"
	"github.com/gokhan/orderly/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(connPool)

	// UseCases
	userUseCase := usecase.NewUserUseCase(store, config)
	productUseCase := usecase.NewProductUseCase(store)
	orderUseCase := usecase.NewOrderUseCase(store)

	// Handlers
	userHandler := http.NewUserHandler(userUseCase)
	productHandler := http.NewProductHandler(productUseCase)
	orderHandler := http.NewOrderHandler(orderUseCase)

	app := fiber.New(fiber.Config{
		AppName: "Orderly API v1.0",
	})

	app.Use(logger.New())
	app.Use(recover.New())

	// Routes
	api := app.Group("/api/v1")
	
	// Public User Routes
	api.Post("/users", userHandler.CreateUser)
	api.Post("/users/login", userHandler.LoginUser)

	// Public Product Routes
	api.Get("/products", productHandler.ListProducts)
	api.Get("/products/:id", productHandler.GetProduct)

	// Protected Routes
	protected := api.Group("/", http.AuthMiddleware(config.TokenSymmetricKey))
	
	// Admin/Protected Product Routes
	protected.Post("/products", productHandler.CreateProduct)

	// Protected Order Routes
	protected.Post("/orders", orderHandler.CreateOrder)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Orderly API is running smoothly",
		})
	})

	log.Printf("Server starting on %s", config.HTTPServerAddress)
	log.Fatal(app.Listen(config.HTTPServerAddress))
}
