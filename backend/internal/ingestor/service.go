package ingestor

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/drewjst/deltagov/internal/congress"
	"github.com/drewjst/deltagov/internal/models"
)

// Service handles bill ingestion from Congress.gov API.
type Service struct {
	db             *gorm.DB
	congressClient *congress.Client
	httpClient     *http.Client
}

// NewService creates a new ingestor service.
func NewService(db *gorm.DB, congressClient *congress.Client) *Service {
	return &Service{
		db:             db,
		congressClient: congressClient,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// IngestResult contains statistics from an ingestion run.
type IngestResult struct {
	BillsFetched   int
	BillsCreated   int
	BillsUpdated   int
	VersionsCreated int
	Errors         []error
}

// IngestRecentBills fetches recent bills from Congress.gov and upserts them.
func (s *Service) IngestRecentBills(ctx context.Context, limit int) (*IngestResult, error) {
	result := &IngestResult{}

	// Fetch recent bills from Congress API
	fetchResult, err := s.congressClient.FetchRecentBills(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("ingestor: failed to fetch recent bills: %w", err)
	}

	result.BillsFetched = len(fetchResult.Bills)
	log.Printf("Fetched %d bills from Congress.gov", result.BillsFetched)

	// Process each bill
	for _, apiBill := range fetchResult.Bills {
		created, updated, versionCreated, err := s.upsertBill(ctx, &apiBill)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("bill %s-%d %s: %w",
				apiBill.Type, apiBill.Congress, apiBill.Number, err))
			continue
		}

		if created {
			result.BillsCreated++
		}
		if updated {
			result.BillsUpdated++
		}
		if versionCreated {
			result.VersionsCreated++
		}
	}

	return result, nil
}

// upsertBill creates or updates a bill and potentially creates a new version.
// Returns (created, updated, versionCreated, error).
func (s *Service) upsertBill(ctx context.Context, apiBill *congress.Bill) (bool, bool, bool, error) {
	// Parse bill number from string
	billNumber, err := strconv.Atoi(apiBill.Number)
	if err != nil {
		return false, false, false, fmt.Errorf("invalid bill number %q: %w", apiBill.Number, err)
	}

	// Convert API bill to metadata JSON
	metadata, err := s.billToMetadata(apiBill)
	if err != nil {
		return false, false, false, fmt.Errorf("failed to create metadata: %w", err)
	}

	// Determine current status from latest action
	currentStatus := ""
	if apiBill.LatestAction != nil {
		currentStatus = apiBill.LatestAction.Text
	}

	// Build the bill model
	bill := models.Bill{
		Congress:       apiBill.Congress,
		BillNumber:     billNumber,
		BillType:       apiBill.Type,
		Title:          apiBill.Title,
		UpdateDate:     apiBill.UpdateDate,
		OriginChamber:  apiBill.OriginChamber,
		CurrentStatus:  currentStatus,
		IsSpendingBill: congress.IsAppropriation(apiBill.Title),
		Metadata:       metadata,
	}

	// Check if bill exists
	var existingBill models.Bill
	err = s.db.WithContext(ctx).
		Where("congress = ? AND bill_number = ? AND bill_type = ?",
			bill.Congress, bill.BillNumber, bill.BillType).
		First(&existingBill).Error

	created := false
	updated := false

	if err == gorm.ErrRecordNotFound {
		// New bill - create it
		if err := s.db.WithContext(ctx).Create(&bill).Error; err != nil {
			return false, false, false, fmt.Errorf("failed to create bill: %w", err)
		}
		created = true
		log.Printf("Created new bill: %s %d (Congress %d)", bill.BillType, bill.BillNumber, bill.Congress)
	} else if err != nil {
		return false, false, false, fmt.Errorf("failed to query bill: %w", err)
	} else {
		// Existing bill - check if UpdateDate changed
		if existingBill.UpdateDate != apiBill.UpdateDate {
			// Update the bill using upsert (ON CONFLICT DO UPDATE)
			bill.ID = existingBill.ID
			if err := s.db.WithContext(ctx).Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "congress"},
					{Name: "bill_number"},
					{Name: "bill_type"},
				},
				DoUpdates: clause.AssignmentColumns([]string{
					"title", "update_date", "origin_chamber",
					"current_status", "is_spending_bill", "metadata", "updated_at",
				}),
			}).Create(&bill).Error; err != nil {
				return false, false, false, fmt.Errorf("failed to update bill: %w", err)
			}
			updated = true
			log.Printf("Updated bill: %s %d (Congress %d) - UpdateDate changed from %s to %s",
				bill.BillType, bill.BillNumber, bill.Congress, existingBill.UpdateDate, apiBill.UpdateDate)
		} else {
			// No changes needed
			bill.ID = existingBill.ID
		}
	}

	// Try to fetch and store bill text as a new version
	versionCreated, err := s.fetchAndStoreVersion(ctx, &bill, apiBill)
	if err != nil {
		// Log but don't fail the entire operation
		log.Printf("Warning: failed to fetch version for %s %d: %v",
			bill.BillType, bill.BillNumber, err)
	}

	return created, updated, versionCreated, nil
}

// fetchAndStoreVersion fetches bill text and creates a version if content is new.
func (s *Service) fetchAndStoreVersion(ctx context.Context, bill *models.Bill, apiBill *congress.Bill) (bool, error) {
	// Parse bill number for API call
	billNumber, _ := strconv.Atoi(apiBill.Number)

	// Fetch text versions from Congress API
	textVersions, err := s.congressClient.GetBillText(ctx, apiBill.Congress, apiBill.Type, billNumber)
	if err != nil {
		// Some bills don't have text yet
		if err == congress.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	if len(textVersions) == 0 {
		return false, nil
	}

	// Get the most recent text version
	latestVersion := textVersions[0]

	// Find a text format URL (prefer XML, then HTML, then TXT)
	textURL := ""
	versionCode := latestVersion.Type
	for _, format := range latestVersion.Formats {
		if format.Type == "Formatted Text" || format.Type == "TXT" {
			textURL = format.URL
			break
		}
		if format.Type == "Formatted XML" || format.Type == "XML" {
			textURL = format.URL
		}
		if textURL == "" && format.Type == "PDF" {
			// Skip PDF for now, can't easily hash
			continue
		}
		if textURL == "" {
			textURL = format.URL
		}
	}

	if textURL == "" {
		return false, nil
	}

	// Fetch the actual text content
	textContent, err := s.fetchTextContent(ctx, textURL)
	if err != nil {
		return false, fmt.Errorf("failed to fetch text from %s: %w", textURL, err)
	}

	// Compute SHA-256 hash
	contentHash := ComputeHash(textContent)

	// Check if we already have this exact version
	var existingVersion models.Version
	err = s.db.WithContext(ctx).
		Where("bill_id = ? AND content_hash = ?", bill.ID, contentHash).
		First(&existingVersion).Error

	if err == nil {
		// Version with same hash already exists
		return false, nil
	} else if err != gorm.ErrRecordNotFound {
		return false, fmt.Errorf("failed to query versions: %w", err)
	}

	// Create new version
	version := models.Version{
		BillID:      bill.ID,
		VersionCode: versionCode,
		ContentHash: contentHash,
		TextContent: textContent,
		FetchedAt:   time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(&version).Error; err != nil {
		return false, fmt.Errorf("failed to create version: %w", err)
	}

	log.Printf("Created new version for %s %d: %s (hash: %s...)",
		bill.BillType, bill.BillNumber, versionCode, contentHash[:16])

	return true, nil
}

// fetchTextContent fetches text content from a URL.
func (s *Service) fetchTextContent(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	// Limit read to 10MB to prevent memory issues
	limited := io.LimitReader(resp.Body, 10*1024*1024)
	content, err := io.ReadAll(limited)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// billToMetadata converts a Congress API bill to a JSONB metadata map.
func (s *Service) billToMetadata(bill *congress.Bill) (datatypes.JSONMap, error) {
	// Marshal to JSON then unmarshal to map for clean conversion
	data, err := json.Marshal(bill)
	if err != nil {
		return nil, err
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return datatypes.JSONMap(metadata), nil
}

// ComputeHash generates a SHA-256 hash of the content.
func ComputeHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}
