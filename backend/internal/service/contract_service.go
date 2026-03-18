package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/broker"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/repository"
	"github.com/google/uuid"
)

// Contrat service defines the business logic for contracts
type ContractService interface {
	CreateContract(ctx context.Context, contract *models.Contract, initialContent string) error
	AssignParticipants(ctx context.Context, contractID uuid.UUID, userID uuid.UUID, participantIDs []uuid.UUID) error
	GetContractDetails(ctx context.Context, contractID uuid.UUID) (*models.Contract, error)
	ListUserContracts(ctx context.Context, userID uuid.UUID) ([]models.Contract, error)
	UpdateContract(ctx context.Context, contractID uuid.UUID, userID uuid.UUID, title string,  description string) error
	DeleteContract(ctx context.Context, contractID uuid.UUID, userId uuid.UUID) error
}

type contractService struct {
    repo      repository.ContractRepository
    publisher *broker.ContractPublisher // <--- Add this line!
}

func NewContractService(repo repository.ContractRepository, pub *broker.ContractPublisher) ContractService {
    return &contractService{
        repo:      repo,
        publisher: pub,
    }
}

// CreateContract create new contract
func (s *contractService) CreateContract(ctx context.Context, contract *models.Contract, initialContent string) error {
	// 1. Logic: Ensure every new contract starts with version 1
	contract.CurrentVersion = 1
	contract.Status = "draft"

	// WRAP the string in a JSON object so Postgres JSONB accepts it
    jsonContent := fmt.Sprintf(`{"content": "%s"}`, initialContent)

	// 2. Initialize the first version snapshot
	firstVersion := models.ContractVersion{
        VersionNumber: 1,
        ContentJSON:   []byte(jsonContent), // Now it's valid JSON
        CreatedBy:     contract.OwnerUserID,
    }
    contract.Versions = append(contract.Versions, firstVersion)
	// 3. Save via Repo
	return s.repo.CreateContract(ctx, contract)
}

func (s *contractService) AssignParticipants(ctx context.Context, contractID uuid.UUID, userID uuid.UUID, participantIDs []uuid.UUID) error {
    // 1. Fetch contract to check ownership
    contract, err := s.repo.GetContractByID(ctx, contractID)
    if err != nil {
        return err
    }

    // 2. Security: Only owner can invite people
    if contract.OwnerUserID != userID {
        return errors.New("unauthorized: you cannot assign participants to this contract")
    }

    // 3. Save to contract_participants table via Repo
    if err := s.repo.AssignParticipants(ctx, contractID, participantIDs); err != nil {
        return err
    }

    // 4. Update status to 'pending' once people are invited
    contract.Status = "pending"
    if err := s.repo.UpdateContract(ctx, contract); err != nil {
        return err
    }

    // 5. COMMUNICATION: Notify Signature Service via RabbitMQ
    event := models.ContractPublishedEvent{
        ContractID:     contractID,
        ParticipantIDs: participantIDs,
        Status:         "pending",
    }
    
    // We only call the publisher once the DB transaction is successful
    return s.publisher.PublishContractPublished(ctx, event)
}

func (s *contractService) GetContractDetails(ctx context.Context, contractID uuid.UUID) (*models.Contract, error) {
	return s.repo.GetContractByID(ctx, contractID)
}

func (s *contractService) ListUserContracts(ctx context.Context, userID uuid.UUID) ([]models.Contract, error) {
	return s.repo.GetContractsByUserID(ctx, userID)
}

// UpdateContract update contract and adds a security check: only the owner can update the contract and for specific criteras
func (s *contractService) UpdateContract(ctx context.Context, contractID uuid.UUID, userID uuid.UUID, title string,  description string) error {
	// 1. Fetch the existing contract to check status and ownership
	contract, err := s.repo.GetContractByID(ctx, contractID)
	if err != nil {
		return err
	}

	// 2. Security Check: Ownership
	if contract.OwnerUserID != userID {
		return  errors.New("unauthorized: you do not own this contract")
	}

	// 3. Business Rule: Lock editing if not in 'draft'
	if contract.Status != "draft" {
		return  errors.New("cannot update a contract that is already pending or completed")
	}

	// 4. Apply Changes
	contract.Title = title
	contract.Description = description

	// 5. Save via Repo
	return s.repo.UpdateContract(ctx, contract)
}

// DeleteContract delete contract and adds a security check: only the owner can delete the contract
func (s *contractService) DeleteContract(ctx context.Context, contractID uuid.UUID, userId uuid.UUID) error {
	// 1. Get the contract
	contract, err := s.repo.GetContractByID(ctx, contractID)
	if err != nil {
		return err
	}

	// 2. Security Check: Only the owner can delete the contract
	if contract.OwnerUserID != userId {
		return errors.New("unauthorized: only the owner can delete the contract")
	}
	return s.repo.DeleteContract(ctx, contractID)
}


