package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtClaims struct {
	jwt.RegisteredClaims
	Claims
}

type jwtManager struct {
	signingMethod jwt.SigningMethod
	key           []byte
	duration      time.Duration
}

func NewJWT(signingMethod jwt.SigningMethod, key []byte, duration time.Duration) Manager {
	return &jwtManager{
		signingMethod: signingMethod,
		key:           key,
		duration:      duration,
	}
}

func (j *jwtManager) CreateToken(userID uuid.UUID, userRole string) (string, error) {
	token, err := jwt.NewWithClaims(j.signingMethod, jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.duration)),
		},
		Claims: Claims{
			UserID:   userID,
			UserRole: userRole,
		},
	}).SignedString(j.key)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (j *jwtManager) VerifyToken(token string) (Claims, error) {
	var cc jwtClaims
	t, err := jwt.ParseWithClaims(token, &cc, func(token *jwt.Token) (any, error) {
		if token.Method != j.signingMethod {
			return nil, fmt.Errorf("%s %w", token.Method.Alg(), jwt.ErrTokenSignatureInvalid)
		}
		return j.key, nil
	})
	if err != nil {
		return Claims{}, err
	}
	if !t.Valid {
		return Claims{}, jwt.ErrTokenNotValidYet
	}
	return cc.Claims, nil
}
