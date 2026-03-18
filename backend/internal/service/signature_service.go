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

type SignatureService interface {
	SignContract(ctx context.Context, sig *models.Signature) error
	GetContractSignatures(ctx context.Context, contractID uuid.UUID) ([]models.Signature, error)
	RevokeSignature(ctx context.Context, signatureID uuid.UUID, userID uuid.UUID) error
}

type signatureService struct {
	repo      repository.SignatureRepository
	publisher *broker.ContractPublisher
}

func NewSignatureService(repo repository.SignatureRepository, pub *broker.ContractPublisher) SignatureService {
	return &signatureService{
		repo:      repo,
		publisher: pub,
	}
}

func (s *signatureService) SignContract(ctx context.Context, sig *models.Signature) error {
	// 1. STATUS CHECK: Ensure contract is 'pending' before allowing a signature
	status, err := s.repo.GetContractStatus(ctx, sig.ContractID)
	if err != nil {
		return errors.New("failed to verify contract status")
	}
	if status != "pending" {
		return fmt.Errorf("cannot sign: contract is %s", status)
	}

	// 2. DATA CHECK: Fixed the string comparison error (sig.FileURL == "")
	if len(sig.VectorData) == 0 && sig.FileURL == "" {
		return errors.New("signature must have either vector data or a file URL")
	}

	// 3. PERSIST: Save to local DB
	if err := s.repo.CreateSignature(ctx, sig); err != nil {
		return err
	}

	// 4. PUBLISH: Notify Contract Service
	event := models.SignatureCreatedEvent{
		SignatureID: sig.ID,
		ContractID:  sig.ContractID,
		UserID:      sig.UserID,
		SignedAt:    sig.SignedAt,
	}

	return s.publisher.PublishSignatureCreated(ctx, event)
}

func (s *signatureService) GetContractSignatures(ctx context.Context, contractID uuid.UUID) ([]models.Signature, error) {
	return s.repo.GetSignatureByContract(ctx, contractID)
}

// 5. REVOKE: Renamed to match interface and added status check
func (s *signatureService) RevokeSignature(ctx context.Context, signatureID uuid.UUID, userID uuid.UUID) error {
	// Fetch signature to find the parent ContractID
	sig, err := s.repo.GetSignatureByID(ctx, signatureID)
	if err != nil {
		return err
	}

	// Status Check: Can't revoke if the contract is already finalized/completed
	status, err := s.repo.GetContractStatus(ctx, sig.ContractID)
	if err != nil {
		return errors.New("failed to verify contract status")
	}
	if status != "pending" {
		return fmt.Errorf("cannot revoke signature: contract is %s", status)
	}

	return s.repo.DeleteSignature(ctx, signatureID, userID)
}