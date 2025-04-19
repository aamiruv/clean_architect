package auth_test

import (
	"testing"
	"time"

	"github.com/amirzayi/clean_architect/pkg/auth"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/require"
)

func TestPaseto(t *testing.T) {
	p := auth.NewPaseto([]byte("YELLOW SUBMARINE, BLACK WIZARDRY"), 1*time.Hour)

	id := uuid.New()
	role := "Admin"

	token, err := p.CreateToken(id, role)
	require.NoError(t, err)

	claims, err := p.VerifyToken(token)
	require.NoError(t, err)

	require.Equal(t, id, claims.UserID)
	require.Equal(t, role, claims.UserRole)
}

func TestPasetoValidation(t *testing.T) {
	p := auth.NewPaseto([]byte("YELLOW SUBMARINE, BLACK WIZARDRY"), -time.Hour)

	id := uuid.New()
	role := "Admin"

	token, err := p.CreateToken(id, role)
	require.NoError(t, err)

	claims, err := p.VerifyToken(token)
	require.ErrorContains(t, err, paseto.ErrTokenValidationError.Error())

	require.Equal(t, uuid.Nil, claims.UserID)
	require.Equal(t, "", claims.UserRole)
}
