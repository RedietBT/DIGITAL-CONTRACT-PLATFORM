package repository

import (
	"context"
	"time"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContractRepository interface {
	CreateContract(ctx context.Context,contract *models.Contract) error
	GetContractByID(ctx context.Context,contractID uuid.UUID) (*models.Contract, error)
	GetContractsByUserID(ctx context.Context,userID uuid.UUID) ([]models.Contract, error)
	UpdateContract(ctx context.Context,contract *models.Contract) error
	DeleteContract(ctx context.Context,contractID uuid.UUID) error
	DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error
	AssignParticipants(ctx context.Context, contractID uuid.UUID, userIDs []uuid.UUID) error
}

type contractRepository struct {
	db *gorm.DB
}

func NewContractRepository(db *gorm.DB) ContractRepository {
	return &contractRepository{db: db}
}

// CreateContract saves the contract AND the first version in one transaction
func (r *contractRepository) CreateContract(ctx context.Context,contract *models.Contract) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. Save the meain contract
		if err := tx.Create(contract).Error; err != nil{
			return err
		}

		// GORM automatically saves associated Versions/Participants 
		// if they are present in the struct!
		return nil
		
	})
}

func (r *contractRepository) GetContractByID(ctx context.Context,contractID uuid.UUID) (*models.Contract, error) {
	var contract models.Contract
	// Preload fetches the versions and participants in the same query
	err := r.db.WithContext(ctx).Preload("Versions").Preload("Participants").First(&contract,"id = ?", contractID).Error

	return &contract, err
}

// GetContractsByUserID retrieves all contracts owned by a specific user
func (r *contractRepository) GetContractsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Contract, error) {
	var contracts []models.Contract
	// We usually don't preload EVERYTHING for a list view to keep it fast, 
	// but we can if you need the participant count.
	err := r.db.WithContext(ctx).Where("owner_user_id = ?", userID).Find(&contracts).Error
	return contracts, err
}

// UpdateContract updates the main contract metadata
func (r *contractRepository) UpdateContract(ctx context.Context, contract *models.Contract) error {
	// Updates() in GORM automatically handles the "updated_at" timestamp
	return r.db.WithContext(ctx).Updates(contract).Error
}

// DeleteContract performs a hard delete of the contract
func (r *contractRepository) DeleteContract(ctx context.Context, contractID uuid.UUID) error {
	// Because of our ON DELETE CASCADE in the database migration,
	// deleting the contract here will automatically clean up versions and participants.
	return r.db.WithContext(ctx).Delete(&models.Contract{}, "id = ?", contractID).Error
}

func (r *contractRepository) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
    // This deletes all contracts where owner_user_id matches the deleted user
    return r.db.WithContext(ctx).Delete(&models.Contract{}, "owner_user_id = ?", userID).Error
}

func (r *contractRepository) AssignParticipants(ctx context.Context, contractID uuid.UUID, participantIDs []uuid.UUID) error {
    // We use a transaction to ensure all participants are added or none are
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        for _, pid := range participantIDs {
            participant := models.ContractParticipant{
                ID:         uuid.New(),
                ContractID: contractID,
                UserID:     pid,
                Role:       "signer", // You can make this dynamic later
                IsRequired:   true,
                AddedAt:    time.Now(),
            }
            
            // Note: 'contract_participants' is the table name from your design doc
            if err := tx.Table("contract_participants").Create(&participant).Error; err != nil {
                return err
            }
        }
        return nil
    })
}