package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type pasetoManager struct {
	key      []byte
	duration time.Duration
}

func NewPaseto(key []byte, duration time.Duration) Manager {
	return &pasetoManager{
		key:      key,
		duration: duration,
	}
}

func (p *pasetoManager) CreateToken(userID uuid.UUID, userRole string) (string, error) {
	jsonToken := paseto.JSONToken{
		Expiration: time.Now().Add(p.duration),
	}
	claims := Claims{
		UserID:   userID,
		UserRole: userRole,
	}
	pasetoMaker := paseto.NewV2()
	token, err := pasetoMaker.Encrypt(p.key, jsonToken, claims)
	return token, err
}

func (p *pasetoManager) VerifyToken(token string) (Claims, error) {
	var (
		jsonToken paseto.JSONToken
		claims    Claims
	)
	pasetoMaker := paseto.NewV2()

	if err := pasetoMaker.Decrypt(token, p.key, &jsonToken, &claims); err != nil {
		return Claims{}, err
	}
	if err := jsonToken.Validate(paseto.ValidAt(time.Now())); err != nil {
		return Claims{}, err
	}

	return claims, nil
}
