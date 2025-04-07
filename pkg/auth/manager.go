package auth

import (
	"github.com/google/uuid"
)

type Claims struct {
	UserID   uuid.UUID `json:"uid"`
	UserRole string    `json:"role"`
}

type Manager interface {
	CreateToken(userID uuid.UUID, userRole string) (token string, err error)
	VerifyToken(token string) (claims Claims, err error)
}
