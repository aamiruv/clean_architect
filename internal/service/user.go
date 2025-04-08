package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/internal/repository"
)

type User interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
}

type user struct {
	db repository.User
}

func NewUserService(db repository.User) User {
	return &user{
		db: db,
	}
}

func (u *user) Create(ctx context.Context, user domain.User) (domain.User, error) {
	user.ID = uuid.New()
	user.Status = domain.UsereStatusNew
	user.CreatedAt = time.Now()

	if err := u.db.Create(ctx, user); err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}
