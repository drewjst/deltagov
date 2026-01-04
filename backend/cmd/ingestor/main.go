package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/drewjst/deltagov/internal/congress"
	"github.com/drewjst/deltagov/internal/database"
	"github.com/drewjst/deltagov/internal/ingestor"
)

func main() {
	// Load .env file if present
	_ = godotenv.Load()

	// Get API key from environment
	apiKey := os.Getenv("CONGRESS_API_KEY")
	if apiKey == "" {
		log.Fatal("CONGRESS_API_KEY environment variable is required")
	}

	// Get database URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Get poll interval (default: 1 hour)
	pollInterval := 1 * time.Hour

	// Connect to database
	dbConfig := database.DefaultConfig(databaseURL)
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)
	log.Println("Connected to database")

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations complete")

	// Create Congress API client
	congressClient, err := congress.New(apiKey)
	if err != nil {
		log.Fatalf("Failed to create Congress client: %v", err)
	}

	// Create ingestor service
	ingestorSvc := ingestor.NewService(db, congressClient)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutdown signal received, stopping ingestor...")
		cancel()
	}()

	log.Println("DeltaGov Ingestor starting...")
	log.Printf("Polling Congress.gov API every %v", pollInterval)

	// Run initial poll
	if err := runIngestion(ctx, ingestorSvc); err != nil {
		log.Printf("Initial ingestion failed: %v", err)
	}

	// Start polling loop
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Ingestor stopped")
			return
		case <-ticker.C:
			if err := runIngestion(ctx, ingestorSvc); err != nil {
				log.Printf("Ingestion failed: %v", err)
			}
		}
	}
}

// runIngestion performs a single ingestion run.
func runIngestion(ctx context.Context, svc *ingestor.Service) error {
	log.Println("Starting ingestion run...")

	result, err := svc.IngestRecentBills(ctx, 50)
	if err != nil {
		return err
	}

	log.Printf("Ingestion complete: fetched=%d, created=%d, updated=%d, versions=%d, errors=%d",
		result.BillsFetched,
		result.BillsCreated,
		result.BillsUpdated,
		result.VersionsCreated,
		len(result.Errors))

	// Log any errors
	for _, e := range result.Errors {
		log.Printf("  Error: %v", e)
	}

	return nil
}
