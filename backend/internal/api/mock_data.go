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

// GetMockHR1 returns mock H.R. 1 data for demo purposes
func GetMockHR1() BillResponse {
	return BillResponse{
		ID:            1,
		Congress:      119,
		BillNumber:    1,
		BillType:      "hr",
		Title:         "One Big Beautiful Bill Act",
		Sponsor:       "Rep. Jason Smith (R-MO)",
		OriginChamber: "House",
		CurrentStatus: "Passed House",
		UpdateDate:    "2025-12-20",
		Versions: []VersionResponse{
			{
				ID:          1,
				VersionCode: "IH",
				Date:        "2025-01-03",
				ContentHash: "abc123",
				Label:       "Introduced in House (Jan 3)",
			},
			{
				ID:          2,
				VersionCode: "RH",
				Date:        "2025-05-15",
				ContentHash: "def456",
				Label:       "Reported in House (May 15)",
			},
			{
				ID:          3,
				VersionCode: "EH",
				Date:        "2025-11-21",
				ContentHash: "ghi789",
				Label:       "Engrossed in House (Nov 21)",
			},
		},
	}
}

// GetMockDiff returns mock diff data for demo purposes
func GetMockDiff() DiffResponse {
	return DiffResponse{
		FromVersion: "IH",
		ToVersion:   "EH",
		Insertions:  156,
		Deletions:   89,
		Lines: []DiffLine{
			{LineNumber: 1, Type: "unchanged", Text: "SECTION 1. SHORT TITLE; TABLE OF CONTENTS."},
			{LineNumber: 2, Type: "unchanged", Text: ""},
			{LineNumber: 3, Type: "unchanged", Text: "(a) SHORT TITLE.—This Act may be cited as the"},
			{LineNumber: 4, Type: "deletion", Text: "\"One Big Beautiful Bill Act of 2025\"."},
			{LineNumber: 5, Type: "insertion", Text: "\"One Big Beautiful Bill Act\"."},
			{LineNumber: 6, Type: "unchanged", Text: ""},
			{LineNumber: 7, Type: "unchanged", Text: "(b) TABLE OF CONTENTS.—The table of contents for this Act is as follows:"},
			{LineNumber: 8, Type: "unchanged", Text: ""},
			{LineNumber: 9, Type: "unchanged", Text: "TITLE I—BORDER SECURITY"},
			{LineNumber: 10, Type: "unchanged", Text: ""},
			{LineNumber: 11, Type: "deletion", Text: "SEC. 101. APPROPRIATIONS FOR BORDER WALL."},
			{LineNumber: 12, Type: "insertion", Text: "SEC. 101. APPROPRIATIONS FOR BORDER SECURITY INFRASTRUCTURE."},
			{LineNumber: 13, Type: "unchanged", Text: ""},
			{LineNumber: 14, Type: "deletion", Text: "There is appropriated $15,000,000,000 for construction of physical"},
			{LineNumber: 15, Type: "deletion", Text: "barriers along the southern border of the United States."},
			{LineNumber: 16, Type: "insertion", Text: "There is appropriated $25,000,000,000 for construction of physical"},
			{LineNumber: 17, Type: "insertion", Text: "barriers, technology systems, and personnel along the southern border"},
			{LineNumber: 18, Type: "insertion", Text: "of the United States."},
			{LineNumber: 19, Type: "unchanged", Text: ""},
			{LineNumber: 20, Type: "unchanged", Text: "SEC. 102. BORDER PATROL AGENTS."},
			{LineNumber: 21, Type: "unchanged", Text: ""},
			{LineNumber: 22, Type: "deletion", Text: "The Secretary of Homeland Security shall hire not fewer than 5,000"},
			{LineNumber: 23, Type: "insertion", Text: "The Secretary of Homeland Security shall hire not fewer than 10,000"},
			{LineNumber: 24, Type: "unchanged", Text: "additional Border Patrol agents within 2 years of enactment."},
			{LineNumber: 25, Type: "unchanged", Text: ""},
			{LineNumber: 26, Type: "unchanged", Text: "TITLE II—TAX PROVISIONS"},
			{LineNumber: 27, Type: "unchanged", Text: ""},
			{LineNumber: 28, Type: "unchanged", Text: "SEC. 201. EXTENSION OF INDIVIDUAL TAX CUTS."},
			{LineNumber: 29, Type: "unchanged", Text: ""},
			{LineNumber: 30, Type: "deletion", Text: "(a) The individual income tax rates established by the Tax Cuts"},
			{LineNumber: 31, Type: "deletion", Text: "and Jobs Act of 2017 are hereby extended through December 31, 2030."},
			{LineNumber: 32, Type: "insertion", Text: "(a) The individual income tax rates established by the Tax Cuts"},
			{LineNumber: 33, Type: "insertion", Text: "and Jobs Act of 2017 are hereby made permanent."},
			{LineNumber: 34, Type: "unchanged", Text: ""},
			{LineNumber: 35, Type: "insertion", Text: "(b) The standard deduction amounts shall be indexed for inflation"},
			{LineNumber: 36, Type: "insertion", Text: "beginning in taxable year 2026."},
			{LineNumber: 37, Type: "unchanged", Text: ""},
			{LineNumber: 38, Type: "unchanged", Text: "SEC. 202. NO TAX ON TIPS."},
			{LineNumber: 39, Type: "unchanged", Text: ""},
			{LineNumber: 40, Type: "unchanged", Text: "Income received as tips by an employee shall not be included in"},
			{LineNumber: 41, Type: "unchanged", Text: "gross income for purposes of the income tax."},
			{LineNumber: 42, Type: "unchanged", Text: ""},
			{LineNumber: 43, Type: "insertion", Text: "SEC. 203. NO TAX ON OVERTIME."},
			{LineNumber: 44, Type: "insertion", Text: ""},
			{LineNumber: 45, Type: "insertion", Text: "Overtime compensation received by an employee shall not be included"},
			{LineNumber: 46, Type: "insertion", Text: "in gross income for purposes of the income tax."},
			{LineNumber: 47, Type: "unchanged", Text: ""},
			{LineNumber: 48, Type: "unchanged", Text: "TITLE III—ENERGY PROVISIONS"},
			{LineNumber: 49, Type: "unchanged", Text: ""},
			{LineNumber: 50, Type: "unchanged", Text: "SEC. 301. REPEAL OF GREEN NEW DEAL PROVISIONS."},
			{LineNumber: 51, Type: "unchanged", Text: ""},
			{LineNumber: 52, Type: "deletion", Text: "(a) The following provisions of the Inflation Reduction Act of 2022"},
			{LineNumber: 53, Type: "deletion", Text: "are hereby repealed: ..."},
			{LineNumber: 54, Type: "insertion", Text: "(a) All climate and clean energy provisions of the Inflation"},
			{LineNumber: 55, Type: "insertion", Text: "Reduction Act of 2022 are hereby repealed effective immediately."},
			{LineNumber: 56, Type: "unchanged", Text: ""},
			{LineNumber: 57, Type: "unchanged", Text: "SEC. 302. DRILLING PERMITS."},
			{LineNumber: 58, Type: "unchanged", Text: ""},
			{LineNumber: 59, Type: "unchanged", Text: "The Secretary of the Interior shall approve all pending applications"},
			{LineNumber: 60, Type: "deletion", Text: "for drilling permits within 60 days of the date of enactment."},
			{LineNumber: 61, Type: "insertion", Text: "for drilling permits within 30 days of the date of enactment."},
		},
		Segments: []DiffSegment{
			{Type: "unchanged", Text: "SECTION 1. SHORT TITLE"},
			{Type: "deletion", Text: "of 2025"},
			{Type: "insertion", Text: ""},
			{Type: "unchanged", Text: "..."},
		},
	}
}
