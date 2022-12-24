package user

import "github.com/google/uuid"

type Authentication struct {
	UserID       uuid.UUID `json:"userID"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	IsFirstLogin bool      `json:"isFirstLogin,omitempty"`
}
