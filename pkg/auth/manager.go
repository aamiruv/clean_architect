package auth

import (
	"github.com/google/uuid"
)

type Claims struct {
	UserID   uuid.UUID
	UserRole string
}

type Manager interface {
	CreateToken(userID uuid.UUID, userRole string) (token string, err error)
	VerifyToken(token string) (claims Claims, err error)
}
