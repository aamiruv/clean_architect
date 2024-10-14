package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/AmirMirzayi/clean_architecture/internal/user"
	"github.com/AmirMirzayi/clean_architecture/internal/user/domain"
	"github.com/google/uuid"
)

func TestCreateUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	u := domain.User{
		Name:        "amir",
		PhoneNumber: "+1234567890",
		Email:       "mirzayi994@gmail.com",
		Password:    "mirzayi994",
	}
	userService := user.NewService(user.NewMemoryRepository())
	result, err := userService.Create(ctx, u)
	if err != nil {
		t.Fail()
	}
	if result.ID == uuid.Nil {
		t.Fail()
	}
	if result.Status != domain.New {
		t.Fail()
	}
	if result.Name != u.Name {
		t.Fail()
	}
	if result.PhoneNumber != u.PhoneNumber {
		t.Fail()
	}
	if result.Password != u.Password {
		t.Fail()
	}
	if result.CreatedAt != time.Now() {
		t.Fail()
	}
	if result.Email != u.Email {
		t.Fail()
	}
}
