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
)

func main() {
	// Load .env file if present
	_ = godotenv.Load()

	// Get API key from environment
	apiKey := os.Getenv("CONGRESS_API_KEY")
	if apiKey == "" {
		log.Fatal("CONGRESS_API_KEY environment variable is required")
	}

	// Get poll interval (default: 1 hour)
	pollInterval := 1 * time.Hour

	// Create Congress API client
	client := congress.NewClient(apiKey)

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
	if err := pollBills(ctx, client); err != nil {
		log.Printf("Initial poll failed: %v", err)
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
			if err := pollBills(ctx, client); err != nil {
				log.Printf("Poll failed: %v", err)
			}
		}
	}
}

// pollBills fetches bills from Congress.gov and stores new versions
func pollBills(ctx context.Context, client *congress.Client) error {
	log.Println("Polling Congress.gov for bill updates...")

	// TODO: Implement actual polling logic
	// 1. Fetch recent bills from Congress API
	// 2. For each bill, compute content hash
	// 3. Compare with stored hash
	// 4. If different, store new version

	bills, err := client.GetRecentBills(ctx)
	if err != nil {
		return err
	}

	log.Printf("Fetched %d bills from Congress.gov", len(bills))
	return nil
}
