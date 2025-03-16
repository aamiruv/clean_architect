package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/amirzayi/clean_architec/internal/user"
	"github.com/amirzayi/clean_architec/internal/user/domain"
)

func TestCreateUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer t.Cleanup(cancel)

	assert := require.New(t)

	u := domain.User{
		Name:        "amir",
		PhoneNumber: "+1234567890",
		Email:       "mirzayi994@gmail.com",
		Password:    "mirzayi994",
	}
	userService := user.NewService(user.NewMemoryRepository())
	result, err := userService.Create(ctx, u)

	assert.Nil(err)
	assert.NotEqual(uuid.Nil, result.ID)
	assert.Equal(domain.New, result.Status)
	assert.Equal(u.Name, result.Name)
	assert.Equal(u.PhoneNumber, result.PhoneNumber)
	assert.Equal(u.Password, result.Password)
	assert.Equal(time.Now(), result.CreatedAt)
	assert.Equal(u.Email, result.Email)
}
