package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Signature struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ContractID    uuid.UUID `gorm:"type:uuid;not null;index" json:"contract_id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	SignatureType string    `gorm:"type:varchar(50);not null" json:"signature_type"` // e.g., 'drawn', 'typed'

	// VectorData stores the coordinates/paths as JSONB
	FileURL string `gorm:"type:text" json:"file_url,omitempty"`
	
	// FIX: Added swaggertype tag for the swagger generator
	VectorData json.RawMessage `gorm:"type:jsonb" json:"vector_data" swaggertype:"object"`
	
	Hash string `gorm:"type:varchar(255)" json:"hash"` // SHA-256 for integrity

	// Position on the PDF page
	PageNumber int     `gorm:"not null;default:1" json:"page_number"`
	PosX       float64 `gorm:"type:numeric(10,2)" json:"pos_x"`
	PosY       float64 `gorm:"type:numeric(10,2)" json:"pos_y"`
	Width      float64 `gorm:"type:numeric(10,2)" json:"width"`
	Height     float64 `gorm:"type:numeric(10,2)" json:"height"`

	RenderStyleID *uuid.UUID `gorm:"type:uuid" json:"render_style_id,omitempty"`

	// Audit Security
	SignedAt   time.Time `gorm:"autoCreateTime" json:"signed_at"`
	IPAddress  string    `gorm:"type:varchar(45)" json:"ip_address"`
	DeviceInfo string    `gorm:"type:text" json:"device_info"`
}

// TableName overrides the default GORM table name to use the schema
func (Signature) TableName() string {
	return "signature_schema.signatures"
}