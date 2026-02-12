package models

import "github.com/google/uuid"

type UserDeletedEvent struct {
	UserID uuid.UUID `json:"user_id"`
	Action string `json:"action"`
}