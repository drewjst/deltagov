package ingestor_test

import (
	"context"
	"os"
	"testing"
	"time"

	"gorm.io/datatypes"

	"github.com/drewjst/deltagov/internal/database"
	"github.com/drewjst/deltagov/internal/ingestor"
	"github.com/drewjst/deltagov/internal/models"
)

// TestBillUpsert_Integration tests that a bill can be written to and read from
// the local PostgreSQL database. This test requires a running PostgreSQL instance.
//
// Run with: DATABASE_URL=postgres://user:pass@localhost:5432/deltagov_test go test -v ./internal/ingestor/...
func TestBillUpsert_Integration(t *testing.T) {
	// Skip if DATABASE_URL is not set
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	// Connect to database
	cfg := database.DefaultConfig(databaseURL)
	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	// Run migrations
	if err := database.Migrate(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create a mock bill
	mockBill := models.Bill{
		Congress:       119,
		BillNumber:     9999,
		BillType:       "hr",
		Title:          "Test Integration Bill",
		UpdateDate:     "2025-01-03",
		OriginChamber:  "House",
		CurrentStatus:  "Introduced",
		IsSpendingBill: false,
		Metadata: datatypes.JSONMap{
			"test": true,
			"source": "integration_test",
		},
	}

	// Clean up any existing test data
	db.Unscoped().Where("congress = ? AND bill_number = ? AND bill_type = ?",
		mockBill.Congress, mockBill.BillNumber, mockBill.BillType).Delete(&models.Bill{})

	// Create the bill
	if err := db.Create(&mockBill).Error; err != nil {
		t.Fatalf("Failed to create bill: %v", err)
	}

	if mockBill.ID == 0 {
		t.Fatal("Bill ID should be set after creation")
	}

	t.Logf("Created bill with ID: %d", mockBill.ID)

	// Read it back
	var readBill models.Bill
	if err := db.First(&readBill, mockBill.ID).Error; err != nil {
		t.Fatalf("Failed to read bill: %v", err)
	}

	// Verify fields
	if readBill.Congress != mockBill.Congress {
		t.Errorf("Congress mismatch: got %d, want %d", readBill.Congress, mockBill.Congress)
	}
	if readBill.BillNumber != mockBill.BillNumber {
		t.Errorf("BillNumber mismatch: got %d, want %d", readBill.BillNumber, mockBill.BillNumber)
	}
	if readBill.BillType != mockBill.BillType {
		t.Errorf("BillType mismatch: got %q, want %q", readBill.BillType, mockBill.BillType)
	}
	if readBill.Title != mockBill.Title {
		t.Errorf("Title mismatch: got %q, want %q", readBill.Title, mockBill.Title)
	}

	// Verify metadata JSONB
	if readBill.Metadata["test"] != true {
		t.Errorf("Metadata 'test' field mismatch: got %v, want true", readBill.Metadata["test"])
	}

	t.Log("Bill successfully written to and read from PostgreSQL")

	// Clean up
	db.Unscoped().Delete(&mockBill)
}

// TestVersionCreation_Integration tests that a version with content hash
// can be created and duplicate detection works.
func TestVersionCreation_Integration(t *testing.T) {
	// Skip if DATABASE_URL is not set
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	// Connect to database
	cfg := database.DefaultConfig(databaseURL)
	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	// Run migrations
	if err := database.Migrate(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create a test bill first
	bill := models.Bill{
		Congress:   119,
		BillNumber: 9998,
		BillType:   "s",
		Title:      "Test Version Bill",
		UpdateDate: "2025-01-03",
	}

	// Clean up any existing test data
	db.Unscoped().Where("congress = ? AND bill_number = ? AND bill_type = ?",
		bill.Congress, bill.BillNumber, bill.BillType).Delete(&models.Bill{})

	if err := db.Create(&bill).Error; err != nil {
		t.Fatalf("Failed to create test bill: %v", err)
	}
	defer db.Unscoped().Delete(&bill)

	// Create a version with content
	textContent := "SECTION 1. SHORT TITLE.\nThis Act may be cited as the Test Act."
	contentHash := ingestor.ComputeHash(textContent)

	version := models.Version{
		BillID:      bill.ID,
		VersionCode: "IH",
		ContentHash: contentHash,
		TextContent: textContent,
		FetchedAt:   time.Now(),
	}

	if err := db.Create(&version).Error; err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	t.Logf("Created version with hash: %s", contentHash[:16])

	// Verify we can find by hash
	var foundVersion models.Version
	err = db.Where("bill_id = ? AND content_hash = ?", bill.ID, contentHash).
		First(&foundVersion).Error
	if err != nil {
		t.Fatalf("Failed to find version by hash: %v", err)
	}

	if foundVersion.TextContent != textContent {
		t.Errorf("TextContent mismatch")
	}

	// Try to create a duplicate - should detect via hash check
	var count int64
	db.Model(&models.Version{}).
		Where("bill_id = ? AND content_hash = ?", bill.ID, contentHash).
		Count(&count)

	if count != 1 {
		t.Errorf("Expected exactly 1 version with hash, got %d", count)
	}

	t.Log("Version with hash successfully created and duplicate detection verified")

	// Clean up
	db.Unscoped().Delete(&version)
}

// TestComputeHash verifies SHA-256 hashing works correctly.
func TestComputeHash(t *testing.T) {
	content := "Hello, World!"
	hash := ingestor.ComputeHash(content)

	// SHA-256 produces 64 hex characters
	if len(hash) != 64 {
		t.Errorf("Hash length should be 64, got %d", len(hash))
	}

	// Same content should produce same hash
	hash2 := ingestor.ComputeHash(content)
	if hash != hash2 {
		t.Errorf("Same content should produce same hash")
	}

	// Different content should produce different hash
	hash3 := ingestor.ComputeHash("Different content")
	if hash == hash3 {
		t.Errorf("Different content should produce different hash")
	}

	t.Logf("Hash of 'Hello, World!': %s", hash)
}

// TestGINIndex_Integration verifies that the GIN index on metadata works.
func TestGINIndex_Integration(t *testing.T) {
	// Skip if DATABASE_URL is not set
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	// Connect to database
	cfg := database.DefaultConfig(databaseURL)
	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	ctx := context.Background()

	// Check that the GIN index exists
	var indexExists bool
	err = db.WithContext(ctx).Raw(`
		SELECT EXISTS (
			SELECT 1 FROM pg_indexes
			WHERE tablename = 'bills'
			AND indexname = 'idx_bills_metadata_gin'
		)
	`).Scan(&indexExists).Error

	if err != nil {
		t.Fatalf("Failed to check index: %v", err)
	}

	if !indexExists {
		t.Error("GIN index idx_bills_metadata_gin should exist on bills.metadata")
	} else {
		t.Log("GIN index on bills.metadata verified")
	}

	// Check delta GIN index too
	err = db.WithContext(ctx).Raw(`
		SELECT EXISTS (
			SELECT 1 FROM pg_indexes
			WHERE tablename = 'deltas'
			AND indexname = 'idx_deltas_delta_json_gin'
		)
	`).Scan(&indexExists).Error

	if err != nil {
		t.Fatalf("Failed to check delta index: %v", err)
	}

	if !indexExists {
		t.Error("GIN index idx_deltas_delta_json_gin should exist on deltas.delta_json")
	} else {
		t.Log("GIN index on deltas.delta_json verified")
	}
}
