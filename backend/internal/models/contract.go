package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Contract represents the master document record
type Contract struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OwnerUserID         uuid.UUID  `gorm:"type:uuid;not null" json:"owner_user_id" validate:"required,min=3,max=100,alphanum_start,no_scripts"`
	Title               string     `gorm:"type:varchar(255);not null" json:"title" validate:"required,min=3"`
	Description         string     `gorm:"type:text" json:"description" validate:"required,min=3,max=1000,no_scripts"`
	Status              string     `gorm:"type:varchar(50);default:'draft'" json:"status"`
	CurrentVersion      int        `grom:"type:int;default:1" json:"current_version"`
	TemplateID          uuid.UUID   `gorm:"type:uuid" json:"template_id,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	CompletedAt         time.Time  `json:"completed_at,omitempty"`

	Versions     []ContractVersion     `gorm:"foreignKey:ContractID" json:"versions,omitempty"`
	Participants []ContractParticipant `gorm:"foreignKey:ContractID" json:"participants,omitempty"`
}

// ContractVersion stores immutable snapshots of the content
type ContractVersion struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ContractID    uuid.UUID      `gorm:"type:uuid;not null" json:"contract_id"`
	VersionNumber int            `gorm:"not null" json:"version_number"`
	ContentJSON   datatypes.JSON `gorm:"type:jsonb;not null" json:"content_json"` // Stores articles/bookmarks
	CreatedBy     uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
}

// ContractParticipant manages roles and signing progress
type ContractParticipant struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ContractID   uuid.UUID  `gorm:"type:uuid;not null" json:"contract_id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"` // Auth Service Reference
	Role         string     `gorm:"type:varchar(50);not null" json:"role"` // signer, viewer, approver
	SigningOrder int        `gorm:"default:0" json:"signing_order"`
	IsRequired   bool       `gorm:"default:true" json:"is_required"`
	Status       string     `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, signed, declined
	SignedAt     *time.Time `json:"signed_at,omitempty"`
	AddedAt      time.Time  `json:"added_at"`
}

// ContractSignatureRequirement defines legal rules for signing
type ContractSignatureRequirement struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ContractID    uuid.UUID `gorm:"type:uuid;not null" json:"contract_id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RequiredTypes []string  `gorm:"type:text[]" json:"required_types"` // handwritten, stamp, initials
	MinRequired   int       `gorm:"default:1" json:"min_required"`
	CreatedAt     time.Time `json:"created_at"`
}

