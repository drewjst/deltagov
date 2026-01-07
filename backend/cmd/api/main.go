package main

import (
	"fmt"
	"log"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/drewjst/deltagov/internal/api"
	"github.com/drewjst/deltagov/internal/congress"
	"github.com/drewjst/deltagov/internal/database"
)

func main() {
	// Load .env file if present
	_ = godotenv.Load()

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Congress client
	congressAPIKey := os.Getenv("CONGRESS_API_KEY")
	var congressClient *congress.Client
	if congressAPIKey != "" {
		var err error
		congressClient, err = congress.NewClient(congress.WithAPIKey(congressAPIKey))
		if err != nil {
			log.Printf("Warning: Failed to create Congress client: %v", err)
		} else {
			log.Println("Congress API client initialized")
		}
	} else {
		log.Println("Warning: CONGRESS_API_KEY not set")
	}

	// Initialize database connection
	var db *gorm.DB
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		dbConfig := database.DefaultConfig(databaseURL)
		var err error
		db, err = database.Connect(dbConfig)
		if err != nil {
			log.Printf("Warning: Failed to connect to database: %v", err)
		} else {
			defer database.Close(db)
			log.Println("Connected to database")

			// Run migrations
			if err := database.Migrate(db); err != nil {
				log.Printf("Warning: Failed to run migrations: %v", err)
			} else {
				log.Println("Database migrations complete")
			}
		}
	} else {
		log.Println("Warning: DATABASE_URL not set, running with mock data only")
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "DeltaGov API",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:4200, http://localhost:80, http://localhost",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Create Huma API with OpenAPI config
	humaConfig := huma.DefaultConfig("DeltaGov API", "1.0.0")
	humaConfig.Info.Description = "API for tracking and comparing legislative bill versions"
	humaConfig.Servers = []*huma.Server{
		{URL: fmt.Sprintf("http://localhost:%s", port), Description: "Local development"},
	}

	humaAPI := humafiber.New(app, humaConfig)

	// Register API routes based on available dependencies
	if db != nil {
		// Database available - register full routes (Congress client optional)
		billService := api.NewBillService(db, congressClient)
		handler := api.NewRouteHandler(billService)
		api.RegisterRoutesWithService(humaAPI, handler)
		log.Println("API routes registered with database support")

		// Register diagnostic routes if Congress client is available
		if congressClient != nil {
			diagnosticSvc := api.NewDiagnosticService(congressClient)
			api.RegisterDiagnosticRoutes(humaAPI, diagnosticSvc)
			log.Println("Diagnostic routes registered")
		}
	} else {
		// Fallback to mock data when no database
		api.RegisterRoutes(humaAPI)
		log.Println("API routes registered with mock data (database not available)")

		// Register diagnostic routes if Congress client is available
		if congressClient != nil {
			diagnosticSvc := api.NewDiagnosticService(congressClient)
			api.RegisterDiagnosticRoutes(humaAPI, diagnosticSvc)
			log.Println("Diagnostic routes registered")
		}
	}

	// Serve Scalar API documentation at /docs
	app.Get("/docs", func(c *fiber.Ctx) error {
		html := `<!DOCTYPE html>
<html>
<head>
    <title>DeltaGov API Docs</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
    <script id="api-reference" data-url="/openapi.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`
		c.Set("Content-Type", "text/html")
		return c.SendString(html)
	})

	// Start server
	log.Printf("DeltaGov API starting on port %s", port)
	log.Printf("API docs available at http://localhost:%s/docs", port)
	log.Printf("OpenAPI spec at http://localhost:%s/openapi.json", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
