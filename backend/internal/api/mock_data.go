package api

import "time"

// MockBill represents bill data for API responses
type MockBill struct {
	ID            string    `json:"id"`
	Congress      int       `json:"congress"`
	Type          string    `json:"type"`
	Number        int       `json:"number"`
	Title         string    `json:"title"`
	Sponsor       string    `json:"sponsor"`
	OriginChamber string    `json:"originChamber"`
	CurrentStatus string    `json:"currentStatus"`
	IntroducedAt  time.Time `json:"introducedAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// MockVersion represents a bill version for API responses
type MockVersion struct {
	ID          string    `json:"id"`
	BillID      string    `json:"billId"`
	Label       string    `json:"label"`
	VersionCode string    `json:"versionCode"`
	ContentHash string    `json:"contentHash"`
	Date        string    `json:"date"`
	FetchedAt   time.Time `json:"fetchedAt"`
}

// MockDiffSegment represents a segment in a diff
type MockDiffSegment struct {
	Type string `json:"type"` // "insertion", "deletion", "unchanged"
	Text string `json:"text"`
	Line int    `json:"line,omitempty"`
}

// MockDelta represents the diff between two versions
type MockDelta struct {
	FromVersion string            `json:"fromVersion"`
	ToVersion   string            `json:"toVersion"`
	Segments    []MockDiffSegment `json:"segments"`
	Insertions  int               `json:"insertions"`
	Deletions   int               `json:"deletions"`
	Unchanged   int               `json:"unchanged"`
}

// GetMockBills returns sample bill data
func GetMockBills() []MockBill {
	return []MockBill{
		{
			ID:            "hr1234-119",
			Congress:      119,
			Type:          "HR",
			Number:        1234,
			Title:         "Federal Budget Act of 2025",
			Sponsor:       "Rep. Jane Smith (D-CA)",
			OriginChamber: "House",
			CurrentStatus: "Passed House",
			IntroducedAt:  time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Date(2025, 12, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:            "s567-119",
			Congress:      119,
			Type:          "S",
			Number:        567,
			Title:         "Infrastructure Investment Act",
			Sponsor:       "Sen. John Doe (R-TX)",
			OriginChamber: "Senate",
			CurrentStatus: "In Committee",
			IntroducedAt:  time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Date(2025, 11, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:            "hr890-119",
			Congress:      119,
			Type:          "HR",
			Number:        890,
			Title:         "Clean Energy Transition Act",
			Sponsor:       "Rep. Maria Garcia (D-NY)",
			OriginChamber: "House",
			CurrentStatus: "Introduced",
			IntroducedAt:  time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC),
			UpdatedAt:     time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC),
		},
	}
}

// GetMockVersions returns sample version data for a bill
func GetMockVersions(billID string) []MockVersion {
	if billID != "hr1234-119" {
		return []MockVersion{}
	}

	return []MockVersion{
		{
			ID:          "v1",
			BillID:      billID,
			Label:       "Version 1 (Dec 1)",
			VersionCode: "IH",
			ContentHash: "abc123def456",
			Date:        "2025-12-01",
			FetchedAt:   time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:          "v2",
			BillID:      billID,
			Label:       "Version 2 (Dec 10)",
			VersionCode: "RH",
			ContentHash: "def456ghi789",
			Date:        "2025-12-10",
			FetchedAt:   time.Date(2025, 12, 10, 14, 30, 0, 0, time.UTC),
		},
		{
			ID:          "v3",
			BillID:      billID,
			Label:       "Version 3 (Dec 15)",
			VersionCode: "EH",
			ContentHash: "ghi789jkl012",
			Date:        "2025-12-15",
			FetchedAt:   time.Date(2025, 12, 15, 9, 0, 0, 0, time.UTC),
		},
		{
			ID:          "v4",
			BillID:      billID,
			Label:       "Version 4 (Dec 20)",
			VersionCode: "ENR",
			ContentHash: "jkl012mno345",
			Date:        "2025-12-20",
			FetchedAt:   time.Date(2025, 12, 20, 16, 0, 0, 0, time.UTC),
		},
	}
}

// GetMockDelta returns a sample diff between two versions
func GetMockDelta(fromVersion, toVersion string) MockDelta {
	return MockDelta{
		FromVersion: fromVersion,
		ToVersion:   toVersion,
		Segments: []MockDiffSegment{
			{Type: "unchanged", Text: "SECTION 1. SHORT TITLE.", Line: 1},
			{Type: "unchanged", Text: "This Act may be cited as the \"Federal Budget Act of 2025\".", Line: 2},
			{Type: "unchanged", Text: "", Line: 3},
			{Type: "unchanged", Text: "SECTION 2. APPROPRIATIONS.", Line: 4},
			{Type: "deletion", Text: "(a) There is appropriated $500,000,000 for infrastructure.", Line: 5},
			{Type: "insertion", Text: "(a) There is appropriated $750,000,000 for infrastructure.", Line: 5},
			{Type: "unchanged", Text: "", Line: 6},
			{Type: "deletion", Text: "(b) Funds shall be distributed over a period of 3 years.", Line: 7},
			{Type: "insertion", Text: "(b) Funds shall be distributed over a period of 5 years.", Line: 7},
			{Type: "unchanged", Text: "", Line: 8},
			{Type: "insertion", Text: "(c) Priority shall be given to rural communities.", Line: 9},
			{Type: "unchanged", Text: "", Line: 10},
			{Type: "unchanged", Text: "SECTION 3. OVERSIGHT.", Line: 11},
			{Type: "unchanged", Text: "The Government Accountability Office shall conduct annual audits.", Line: 12},
		},
		Insertions: 3,
		Deletions:  2,
		Unchanged:  9,
	}
}
