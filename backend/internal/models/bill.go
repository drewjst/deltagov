package models

import (
	"time"

	"gorm.io/datatypes"
)

// Bill represents a legislative bill with GORM ORM mappings.
// The composite unique key is (Congress, BillNumber, BillType).
type Bill struct {
	ID             uint              `json:"id" gorm:"primaryKey"`
	Congress       int               `json:"congress" gorm:"uniqueIndex:idx_bill_unique,priority:1"`
	BillNumber     int               `json:"bill_number" gorm:"uniqueIndex:idx_bill_unique,priority:2"`
	BillType       string            `json:"bill_type" gorm:"uniqueIndex:idx_bill_unique,priority:3;size:10"`
	Title          string            `json:"title"`
	Sponsor        string            `json:"sponsor,omitempty"`
	OriginChamber  string            `json:"origin_chamber"`
	CurrentStatus  string            `json:"current_status"`
	UpdateDate     string            `json:"update_date"` // Congress.gov updateDate string
	IsSpendingBill bool              `json:"is_spending_bill" gorm:"index"`
	Metadata       datatypes.JSONMap `json:"metadata" gorm:"type:jsonb"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// Version represents a point-in-time snapshot of bill text.
// Uses SHA-256 content hash for deduplication.
type Version struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	BillID      uint      `json:"bill_id" gorm:"index"`
	VersionCode string    `json:"version_code"` // e.g., "IH" (Introduced House), "EH" (Engrossed House)
	ContentHash string    `json:"content_hash" gorm:"index;size:64"` // SHA-256 hash
	TextContent string    `json:"text_content" gorm:"type:text"`
	FetchedAt   time.Time `json:"fetched_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// Delta represents a stored diff between two versions.
// DeltaJSON stores structured diff data as JSONB for querying.
type Delta struct {
	ID         uint              `json:"id" gorm:"primaryKey"`
	VersionAID uint              `json:"version_a_id" gorm:"index"`
	VersionBID uint              `json:"version_b_id" gorm:"index"`
	Insertions int               `json:"insertions"`
	Deletions  int               `json:"deletions"`
	DeltaJSON  datatypes.JSONMap `json:"delta_json" gorm:"type:jsonb"` // Structured diff data
	ComputedAt time.Time         `json:"computed_at"`
	CreatedAt  time.Time         `json:"created_at"`
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
