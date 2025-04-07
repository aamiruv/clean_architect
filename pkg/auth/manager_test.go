package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewJWT(t *testing.T) {
	m := NewJWT("sample_key", time.Hour)

	id := uuid.New()
	role := "Admin"
	token, err := m.CreateToken(id, role)
	require.NoError(t, err)

	claims, err := m.VerifyToken(token)
	require.NoError(t, err)

	require.Equal(t, id, claims.UserID)
	require.Equal(t, role, claims.UserRole)
}
