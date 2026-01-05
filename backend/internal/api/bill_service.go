package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/drewjst/deltagov/internal/congress"
	"github.com/drewjst/deltagov/internal/diff_engine"
	"github.com/drewjst/deltagov/internal/models"
	"gorm.io/gorm"
)

// BillService handles bill-related business logic.
type BillService struct {
	db             *gorm.DB
	congressClient *congress.Client
}

// NewBillService creates a new BillService instance.
func NewBillService(db *gorm.DB, congressClient *congress.Client) *BillService {
	return &BillService{
		db:             db,
		congressClient: congressClient,
	}
}

// BillResponse is the API response format for a bill.
type BillResponse struct {
	ID            uint              `json:"id"`
	Congress      int               `json:"congress"`
	BillNumber    int               `json:"billNumber"`
	BillType      string            `json:"billType"`
	Title         string            `json:"title"`
	Sponsor       string            `json:"sponsor"`
	OriginChamber string            `json:"originChamber"`
	CurrentStatus string            `json:"currentStatus"`
	UpdateDate    string            `json:"updateDate"`
	Versions      []VersionResponse `json:"versions,omitempty"`
}

// VersionResponse is the API response format for a version.
type VersionResponse struct {
	ID          uint   `json:"id"`
	VersionCode string `json:"versionCode"`
	Date        string `json:"date"`
	ContentHash string `json:"contentHash"`
	Label       string `json:"label"`
}

// DiffResponse is the API response format for a diff.
type DiffResponse struct {
	FromVersion string        `json:"fromVersion"`
	ToVersion   string        `json:"toVersion"`
	Insertions  int           `json:"insertions"`
	Deletions   int           `json:"deletions"`
	Lines       []DiffLine    `json:"lines"`
	Segments    []DiffSegment `json:"segments"`
}

// DiffLine represents a single line in the diff output.
type DiffLine struct {
	LineNumber int    `json:"lineNumber"`
	Type       string `json:"type"` // "insertion", "deletion", "unchanged"
	Text       string `json:"text"`
}

// DiffSegment represents a segment in the diff output (word-level).
type DiffSegment struct {
	Type string `json:"type"` // "insertion", "deletion", "unchanged"
	Text string `json:"text"`
}

// versionCodeLabels maps version codes to human-readable labels.
var versionCodeLabels = map[string]string{
	"IH":  "Introduced in House",
	"RH":  "Reported in House",
	"EH":  "Engrossed in House",
	"IS":  "Introduced in Senate",
	"RS":  "Reported in Senate",
	"ES":  "Engrossed in Senate",
	"PCS": "Placed on Calendar Senate",
	"EAS": "Engrossed Amendment Senate",
	"ENR": "Enrolled",
	"PL":  "Public Law",
}

// FetchAndStoreHR1 fetches H.R. 1 (119th Congress) and stores it in the database.
// This is the "One Big Beautiful Bill".
func (s *BillService) FetchAndStoreHR1(ctx context.Context) (*BillResponse, error) {
	const (
		congressNum = 119
		billType    = "hr"
		billNumber  = 1
	)

	// Check if we already have this bill in the database
	var existingBill models.Bill
	result := s.db.Where("congress = ? AND bill_type = ? AND bill_number = ?",
		congressNum, billType, billNumber).First(&existingBill)

	if result.Error == nil {
		// Bill exists, check if we need to refresh versions
		var versionCount int64
		s.db.Model(&models.Version{}).Where("bill_id = ?", existingBill.ID).Count(&versionCount)

		if versionCount > 0 {
			// Return existing bill with versions
			return s.GetBillWithVersions(ctx, existingBill.ID)
		}
	}

	// Fetch bill details from Congress.gov
	log.Printf("Fetching H.R. 1 (119th Congress) from Congress.gov...")
	billDetail, err := s.congressClient.GetBillDetail(ctx, congressNum, billType, billNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bill details: %w", err)
	}

	// Create or update the bill record
	bill := models.Bill{
		Congress:      congressNum,
		BillNumber:    billNumber,
		BillType:      billType,
		Title:         billDetail.Title,
		OriginChamber: billDetail.OriginChamber,
		UpdateDate:    billDetail.UpdateDate,
	}

	if billDetail.LatestAction != nil {
		bill.CurrentStatus = billDetail.LatestAction.Text
	}

	// Upsert the bill
	if result.Error != nil {
		// Create new bill
		if err := s.db.Create(&bill).Error; err != nil {
			return nil, fmt.Errorf("failed to create bill: %w", err)
		}
	} else {
		// Update existing bill
		bill.ID = existingBill.ID
		if err := s.db.Save(&bill).Error; err != nil {
			return nil, fmt.Errorf("failed to update bill: %w", err)
		}
	}

	// Fetch all text versions with content
	log.Printf("Fetching text versions for H.R. 1...")
	textVersions, err := s.congressClient.GetBillTextWithContent(ctx, congressNum, billType, billNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch text versions: %w", err)
	}

	log.Printf("Found %d text versions", len(textVersions))

	// Store each version
	for _, tv := range textVersions {
		// Compute content hash
		hash := sha256.Sum256([]byte(tv.Content))
		contentHash := hex.EncodeToString(hash[:])

		// Extract version code from type (e.g., "Introduced in House" -> "IH")
		versionCode := extractVersionCode(tv.Type)

		// Check if version already exists
		var existingVersion models.Version
		if err := s.db.Where("bill_id = ? AND version_code = ?", bill.ID, versionCode).
			First(&existingVersion).Error; err == nil {
			// Version exists, skip
			continue
		}

		// Parse date
		fetchedAt := time.Now()
		if tv.Date != "" {
			if parsed, err := time.Parse("2006-01-02", tv.Date); err == nil {
				fetchedAt = parsed
			}
		}

		version := models.Version{
			BillID:      bill.ID,
			VersionCode: versionCode,
			ContentHash: contentHash,
			TextContent: tv.Content,
			FetchedAt:   fetchedAt,
		}

		if err := s.db.Create(&version).Error; err != nil {
			log.Printf("Warning: failed to create version %s: %v", versionCode, err)
			continue
		}
		log.Printf("Stored version: %s (%s)", versionCode, tv.Type)
	}

	return s.GetBillWithVersions(ctx, bill.ID)
}

// GetBillWithVersions retrieves a bill with all its versions.
func (s *BillService) GetBillWithVersions(ctx context.Context, billID uint) (*BillResponse, error) {
	var bill models.Bill
	if err := s.db.First(&bill, billID).Error; err != nil {
		return nil, fmt.Errorf("bill not found: %w", err)
	}

	var versions []models.Version
	// Select specific fields to avoid fetching large text_content
	if err := s.db.Select("id", "bill_id", "version_code", "content_hash", "fetched_at").
		Where("bill_id = ?", billID).Order("fetched_at ASC").Find(&versions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch versions: %w", err)
	}

	response := &BillResponse{
		ID:            bill.ID,
		Congress:      bill.Congress,
		BillNumber:    bill.BillNumber,
		BillType:      bill.BillType,
		Title:         bill.Title,
		Sponsor:       bill.Sponsor,
		OriginChamber: bill.OriginChamber,
		CurrentStatus: bill.CurrentStatus,
		UpdateDate:    bill.UpdateDate,
		Versions:      make([]VersionResponse, len(versions)),
	}

	for i, v := range versions {
		label := versionCodeLabels[v.VersionCode]
		if label == "" {
			label = v.VersionCode
		}
		response.Versions[i] = VersionResponse{
			ID:          v.ID,
			VersionCode: v.VersionCode,
			Date:        v.FetchedAt.Format("2006-01-02"),
			ContentHash: v.ContentHash,
			Label:       fmt.Sprintf("%s (%s)", label, v.FetchedAt.Format("Jan 2")),
		}
	}

	return response, nil
}

// ComputeDiff computes a diff between two versions.
func (s *BillService) ComputeDiff(ctx context.Context, fromVersionID, toVersionID uint) (*DiffResponse, error) {
	var fromVersion, toVersion models.Version

	if err := s.db.First(&fromVersion, fromVersionID).Error; err != nil {
		return nil, fmt.Errorf("from version not found: %w", err)
	}
	if err := s.db.First(&toVersion, toVersionID).Error; err != nil {
		return nil, fmt.Errorf("to version not found: %w", err)
	}

	// Check if we have a cached delta
	var existingDelta models.Delta
	if err := s.db.Where("version_a_id = ? AND version_b_id = ?",
		fromVersionID, toVersionID).First(&existingDelta).Error; err == nil {
		// Return cached delta
		return s.deltaToResponse(&existingDelta, fromVersion.VersionCode, toVersion.VersionCode), nil
	}

	// Compute the diff using the diff engine
	delta, err := diff_engine.ComputeWordLevel(fromVersion.TextContent, toVersion.TextContent)
	if err != nil {
		return nil, fmt.Errorf("failed to compute diff: %w", err)
	}

	// Store the delta for caching
	storedDelta := models.Delta{
		VersionAID: fromVersionID,
		VersionBID: toVersionID,
		Insertions: delta.Insertions,
		Deletions:  delta.Deletions,
		ComputedAt: time.Now(),
	}
	s.db.Create(&storedDelta)

	// Convert to response format
	response := &DiffResponse{
		FromVersion: fromVersion.VersionCode,
		ToVersion:   toVersion.VersionCode,
		Insertions:  delta.Insertions,
		Deletions:   delta.Deletions,
		Lines:       make([]DiffLine, 0, len(delta.Hunks)*10),
		Segments:    make([]DiffSegment, 0),
	}

	lineNum := 1
	for _, hunk := range delta.Hunks {
		for _, change := range hunk.Lines {
			changeType := "unchanged"
			switch change.Type {
			case diff_engine.ChangeInsert:
				changeType = "insertion"
			case diff_engine.ChangeDelete:
				changeType = "deletion"
			case diff_engine.ChangeUnchanged:
				changeType = "unchanged"
			}

			response.Lines = append(response.Lines, DiffLine{
				LineNumber: lineNum,
				Type:       changeType,
				Text:       change.Content,
			})
			response.Segments = append(response.Segments, DiffSegment{
				Type: changeType,
				Text: change.Content,
			})
			lineNum++
		}
	}

	return response, nil
}

// deltaToResponse converts a stored Delta to DiffResponse.
func (s *BillService) deltaToResponse(delta *models.Delta, fromCode, toCode string) *DiffResponse {
	return &DiffResponse{
		FromVersion: fromCode,
		ToVersion:   toCode,
		Insertions:  delta.Insertions,
		Deletions:   delta.Deletions,
		Lines:       []DiffLine{},
		Segments:    []DiffSegment{},
	}
}

// extractVersionCode extracts the version code from the full type string.
func extractVersionCode(typeStr string) string {
	// Map full type names to codes
	typeToCode := map[string]string{
		"Introduced in House":       "IH",
		"Reported in House":         "RH",
		"Engrossed in House":        "EH",
		"Introduced in Senate":      "IS",
		"Reported in Senate":        "RS",
		"Engrossed in Senate":       "ES",
		"Placed on Calendar Senate": "PCS",
		"Engrossed Amendment Senate": "EAS",
		"Enrolled":                  "ENR",
		"Public Law":                "PL",
	}

	if code, ok := typeToCode[typeStr]; ok {
		return code
	}

	// If not found, return first two letters uppercase
	if len(typeStr) >= 2 {
		return typeStr[:2]
	}
	return typeStr
}

// GetAllBills returns all bills from the database.
func (s *BillService) GetAllBills(ctx context.Context) ([]BillResponse, error) {
	var bills []models.Bill
	if err := s.db.Find(&bills).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch bills: %w", err)
	}

	responses := make([]BillResponse, len(bills))
	for i, b := range bills {
		responses[i] = BillResponse{
			ID:            b.ID,
			Congress:      b.Congress,
			BillNumber:    b.BillNumber,
			BillType:      b.BillType,
			Title:         b.Title,
			Sponsor:       b.Sponsor,
			OriginChamber: b.OriginChamber,
			CurrentStatus: b.CurrentStatus,
			UpdateDate:    b.UpdateDate,
		}
	}

	return responses, nil
}

// GetBillByID retrieves a single bill by its database ID.
func (s *BillService) GetBillByID(ctx context.Context, id uint) (*BillResponse, error) {
	return s.GetBillWithVersions(ctx, id)
}
