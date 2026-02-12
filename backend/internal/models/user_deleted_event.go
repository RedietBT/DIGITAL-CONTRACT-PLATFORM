package models

import "github.com/google/uuid"

type UserDeletedEvent struct {
	UserId uuid.UUID `json:"user_id"`
	Action string `json:"action"`
}