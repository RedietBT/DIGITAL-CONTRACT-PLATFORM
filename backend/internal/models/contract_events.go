package models

import "github.com/google/uuid"

type ContractPublishedEvent struct {
	ContractID uuid.UUID `json:"contract_id"`
	ParticipantIDs []uuid.UUID `json:"participant_ids"`
	Status string `json:"status"`
}