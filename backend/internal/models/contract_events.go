package models

import (
	"time"

	"github.com/google/uuid"
)

type ContractPublishedEvent struct {
	ContractID uuid.UUID `json:"contract_id"`
	ParticipantIDs []uuid.UUID `json:"participant_ids"`
	Status string `json:"status"`
}

// New Event (Signature -> Contract)
type SignatureCreatedEvent struct {
	SignatureID uuid.UUID `json:"signature_id"`
	ContractID  uuid.UUID `json:"contract_id"`
	UserID      uuid.UUID `json:"user_id"`
	SignedAt    time.Time `json:"signed_at"`
}