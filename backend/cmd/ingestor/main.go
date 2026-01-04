package main

import (
	"context"
	"flag"
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
	// Parse command-line flags
	singleRun := flag.Bool("single-run", false, "Run ingestion once and exit (for Cloud Run Jobs)")
	billLimit := flag.Int("limit", 50, "Maximum number of bills to fetch per run")
	flag.Parse()

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

	// Get poll interval from environment (default: 1 hour)
	pollInterval := 1 * time.Hour
	if intervalStr := os.Getenv("POLL_INTERVAL"); intervalStr != "" {
		if parsed, err := time.ParseDuration(intervalStr); err == nil {
			pollInterval = parsed
		}
	}

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

	// Single-run mode for Cloud Run Jobs
	if *singleRun {
		log.Println("DeltaGov Ingestor running in single-run mode...")
		if err := runIngestion(ctx, ingestorSvc, *billLimit); err != nil {
			log.Fatalf("Ingestion failed: %v", err)
		}
		log.Println("Single-run ingestion complete, exiting")
		return
	}

	// Continuous polling mode
	log.Println("DeltaGov Ingestor starting in continuous mode...")
	log.Printf("Polling Congress.gov API every %v", pollInterval)

	// Run initial poll
	if err := runIngestion(ctx, ingestorSvc, *billLimit); err != nil {
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
			if err := runIngestion(ctx, ingestorSvc, *billLimit); err != nil {
				log.Printf("Ingestion failed: %v", err)
			}
		}
	}
}

// runIngestion performs a single ingestion run.
func runIngestion(ctx context.Context, svc *ingestor.Service, limit int) error {
	log.Printf("Starting ingestion run (limit=%d)...", limit)

	result, err := svc.IngestRecentBills(ctx, limit)
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
