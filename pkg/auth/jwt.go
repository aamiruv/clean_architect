package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtClaims struct {
	jwt.RegisteredClaims
	Claims
}

type jwtManager struct {
	key      string
	duration time.Duration
}

func NewJWT(key string, duration time.Duration) Manager {
	return &jwtManager{
		key:      key,
		duration: duration,
	}
}

func (j *jwtManager) CreateToken(userID uuid.UUID, userRole string) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.duration)),
		},
		Claims: Claims{
			UserID:   userID,
			UserRole: userRole,
		},
	}).SignedString([]byte(j.key))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (j *jwtManager) VerifyToken(token string) (Claims, error) {
	var cc jwtClaims
	t, err := jwt.ParseWithClaims(token, &cc, func(t *jwt.Token) (any, error) {
		return []byte(j.key), nil
	})
	if err != nil {
		return Claims{}, err
	}
	if !t.Valid {
		return Claims{}, jwt.ErrTokenNotValidYet
	}
	return cc.Claims, nil
}
