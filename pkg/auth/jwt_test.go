package auth_test

import (
	"testing"
	"time"

	"github.com/amirzayi/clean_architect/pkg/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	m := auth.NewJWT(jwt.SigningMethodHS512, []byte("sample_key"), time.Hour)

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
	m := auth.NewJWT(jwt.SigningMethodHS512, []byte("sample_key"), -time.Hour)

	id := uuid.New()
	role := "Admin"
	token, err := m.CreateToken(id, role)
	require.NoError(t, err)

	claims, err := m.VerifyToken(token)
	require.ErrorIs(t, err, jwt.ErrTokenExpired)

	require.Equal(t, uuid.Nil, claims.UserID)
	require.Equal(t, "", claims.UserRole)
}
func TestJWTValidation_DifferentSignMethod(t *testing.T) {
	m := auth.NewJWT(jwt.SigningMethodHS512, []byte("sample_key"), time.Hour)

	id := uuid.New()
	role := "Admin"
	token, err := m.CreateToken(id, role)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	newM := auth.NewJWT(jwt.SigningMethodHS256, []byte("sample_key"), time.Hour)
	newToken, err := newM.CreateToken(id, role)
	require.NoError(t, err)

	claims, err := m.VerifyToken(newToken)
	require.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)

	require.Empty(t, claims)
}
