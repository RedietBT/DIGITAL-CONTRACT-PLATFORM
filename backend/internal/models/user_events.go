package models

//UserCreatedEvent is the data sent when a new user registers
type UserCreatedEvent struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	DisplayName string `json:"display_name"`
}