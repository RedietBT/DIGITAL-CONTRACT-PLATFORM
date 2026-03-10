package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)


type Signature struct {
	ID      		uuid.UUID    `gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	ContractID      uuid.UUID    `gorm:"type:uuid;not null;index"`
	UserID          uuid.UUID    `gorm:"type:uuid;not null;index"`
	SignatureType   string       `gorm:"type:varchar(50);not null"` // e.g., 'drawn', 'typed'

	// VectorData stores the coordinates/paths as JSONB
	FileURL         string        `gorm:"type:text"`
	VectorData      json.RawMessage `gorm:"type:jsonb"`
	Hash            string          `gorm:"type:varchar(255)"` // SHA-256 for integrity

	// Position on the PDF page
	PageNumber      int             `gorm:"not null;default:1"`
	PosX            float64         `gorm:"type:numeric(10,2)"`
	PosY            float64         `gorm:"type:numeric(10,2)"`
	Width           float64         `gorm:"type:numeric(10,2)"`
	Height          float64         `gorm:"type:numeric(10,2)"`

	RenderStyleID   *uuid.UUID      `gorm:"type:uuid"`


	// Audit Security
	SignedAt        time.Time       `gorm:"autoCreateTime"`
	IPAddress       string          `gorm:"type:varchar(45)"`
	DeviceInfo      string          `gorm:"type:text"`
}

// TableName overrides the default GORM table name to use the schema
func (Signature) TableName() string {
	return "signature_schema.signatures"
}
