package user

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Authentication contains details which would to sent as response on login.
type Authentication struct {
	UserID       uuid.UUID `json:"userID"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	IsFirstLogin bool      `json:"isFirstLogin,omitempty"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	UserID uuid.UUID `json:"userID"`
	jwt.RegisteredClaims
}
