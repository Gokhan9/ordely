package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	_ "github.com/gokhan/orderly/docs"
	"github.com/gokhan/orderly/internal/delivery/http"
	"github.com/gokhan/orderly/internal/repository/db"
	"github.com/gokhan/orderly/internal/usecase"
	"github.com/gokhan/orderly/pkg/config"
	pkgLogger "github.com/gokhan/orderly/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// @title Orderly API
// @version 1.0
// @description High performance e-commerce API.
// @termsOfService http://swagger.io/terms/

// @contact.name Gökhan
// @contact.email gokhan@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and then your token.

func main() {
	pkgLogger.SetupLogger()

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}
	defer connPool.Close()

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

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Orderly API is running smoothly",
		})
	})

	// Graceful Shutdown
	go func() {
		log.Info().Msgf("Server starting on %s", config.HTTPServerAddress)
		if err := app.Listen(config.HTTPServerAddress); err != nil {
			log.Fatal().Err(err).Msg("server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("Server exited properly")
}
