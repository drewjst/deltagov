package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"github.com/drewjst/deltagov/internal/api"
)

func main() {
	// Load .env file if present
	_ = godotenv.Load()

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "DeltaGov API",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:4200",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "deltagov-api",
		})
	})

	// Register API routes
	api.RegisterRoutes(app)

	// Start server
	log.Printf("DeltaGov API starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
