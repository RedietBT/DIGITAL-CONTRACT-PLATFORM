package repository

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SignatureRepository interface {
	CreateSignature(ctx context.Context, sig *models.Signature) error
	GetSignatureByID(ctx context.Context, id uuid.UUID) (*models.Signature, error)
	GetSignatureByContract(ctx context.Context, contractID uuid.UUID) ([]models.Signature, error)
	DeleteSignature(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetContractStatus(ctx context.Context, contractID uuid.UUID) (string, error)
}

type signatureRepository struct {
	db *gorm.DB
}

func NewSignatureRepository(db *gorm.DB) SignatureRepository {
	return &signatureRepository{db: db}
}

func (r *signatureRepository) CreateSignature(ctx context.Context, sig *models.Signature) error {
	// 1. Generate Hash of the VectorData for integrity
	if len(sig.VectorData) > 0 {
		hash := sha256.Sum256(sig.VectorData)
		sig.Hash = fmt.Sprintf("%x", hash)
	}

	// 2. Persist to DB
	return  r.db.WithContext(ctx).Create(sig).Error
}

// GetSignatureByID fetches a specific signature mark
func (r *signatureRepository) GetSignatureByID(ctx context.Context, id uuid.UUID) (*models.Signature, error) {
	var sig models.Signature
	err := r.db.WithContext(ctx).First(&sig, "id = ?", id).Error
	return &sig, err
}


// GetSignaturesByContract fetches all signatures for a document (to render them all)
func (r *signatureRepository) GetSignatureByContract(ctx context.Context, contractID uuid.UUID) ([]models.Signature, error){
	var signatures []models.Signature
	err := r.db.WithContext(ctx).Where("contract_id = ?", contractID).Find(&signatures).Error
	return signatures, err
}

// DeleteSignature allows a user to remove their mark BEFORE the contract is finalized
func (r *signatureRepository) DeleteSignature(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// We include userID to ensure only the owner can delete their own signature
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.Signature{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *signatureRepository) GetContractStatus(ctx context.Context, contractID uuid.UUID) (string, error) {
    var status string
    // This assumes you have a table 'contract_permissions' or 'contracts' in this schema
    err := r.db.WithContext(ctx).Table("contract_permissions").
        Select("status").Where("contract_id = ?", contractID).
        Scan(&status).Error
    return status, err
}


