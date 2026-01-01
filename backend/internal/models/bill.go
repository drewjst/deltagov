package models

import (
	"time"
)

// Bill represents a legislative bill
type Bill struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	Congress      int       `json:"congress"`
	Type          string    `json:"type"`
	Number        int       `json:"number"`
	Title         string    `json:"title"`
	Sponsor       string    `json:"sponsor,omitempty"`
	OriginChamber string    `json:"origin_chamber"`
	CurrentStatus string    `json:"current_status"`
	IntroducedAt  time.Time `json:"introduced_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// Version represents a point-in-time snapshot of bill text
type Version struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	BillID      string    `json:"bill_id" gorm:"index"`
	VersionCode string    `json:"version_code"` // e.g., "IH" (Introduced House), "EH" (Engrossed House)
	ContentHash string    `json:"content_hash"` // SHA-256 hash
	TextContent string    `json:"text_content" gorm:"type:text"`
	FetchedAt   time.Time `json:"fetched_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// Delta represents a stored diff between two versions
type Delta struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	VersionAID string    `json:"version_a_id" gorm:"index"`
	VersionBID string    `json:"version_b_id" gorm:"index"`
	Insertions int       `json:"insertions"`
	Deletions  int       `json:"deletions"`
	DeltaJSON  string    `json:"delta_json" gorm:"type:jsonb"` // Structured diff data
	ComputedAt time.Time `json:"computed_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName returns the table name for Bill
func (Bill) TableName() string {
	return "bills"
}

// TableName returns the table name for Version
func (Version) TableName() string {
	return "versions"
}

// TableName returns the table name for Delta
func (Delta) TableName() string {
	return "deltas"
}
