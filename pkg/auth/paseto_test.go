package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPaseto(t *testing.T) {
	p := NewPaseto([]byte("YELLOW SUBMARINE, BLACK WIZARDRY"), 1*time.Hour)

	id := uuid.New()
	role := "Admin"

	token, err := p.CreateToken(id, role)
	require.NoError(t, err)

	claims, err := p.VerifyToken(token)
	require.NoError(t, err)

	require.Equal(t, id, claims.UserID)
	require.Equal(t, role, claims.UserRole)
}
