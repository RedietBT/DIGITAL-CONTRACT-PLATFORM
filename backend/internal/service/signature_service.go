package service

import (
	"context"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/repository"
	"github.com/google/uuid"
)

type SignatureService interface {
	SignDocument(ctx context.Context, sig *models.Signature) error
	GetContractSignatures(ctx context.Context, contractID uuid.UUID) ([]models.Signature, error)
	RevokeSignature(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type signatureService struct {
	repo repository.SignatureRepository
}

// GetContractSignatures implements [SignatureService].
func (s *signatureService) GetContractSignatures(ctx context.Context, contractID uuid.UUID) ([]models.Signature, error) {
	panic("unimplemented")
}

// RevokeSignature implements [SignatureService].
func (s *signatureService) RevokeSignature(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	panic("unimplemented")
}

func NewSignatureService(repo repository.SignatureRepository) SignatureService {
	return &signatureService{repo: repo}
}

func (s *signatureService) SignDocument(ctx context.Context, sig *models.Signature) error {

} 
