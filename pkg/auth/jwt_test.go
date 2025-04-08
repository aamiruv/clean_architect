package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	m := NewJWT(jwt.SigningMethodHS512, []byte("sample_key"), time.Hour)

	id := uuid.New()
	role := "Admin"
	token, err := m.CreateToken(id, role)
	require.NoError(t, err)

	claims, err := m.VerifyToken(token)
	require.NoError(t, err)

	require.Equal(t, id, claims.UserID)
	require.Equal(t, role, claims.UserRole)
}

func TestJWTValidation(t *testing.T) {
	m := NewJWT(jwt.SigningMethodHS512, []byte("sample_key"), -time.Hour)

	id := uuid.New()
	role := "Admin"
	token, err := m.CreateToken(id, role)
	require.NoError(t, err)

	claims, err := m.VerifyToken(token)
	require.ErrorIs(t, err, jwt.ErrTokenExpired)

	require.Equal(t, uuid.Nil, claims.UserID)
	require.Equal(t, "", claims.UserRole)
}
