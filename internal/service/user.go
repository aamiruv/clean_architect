package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/internal/repository"
	"github.com/amirzayi/clean_architect/pkg/cache"
)

type User interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}

type user struct {
	db    repository.User
	cache cache.Cache[domain.User]
}

func NewUserService(db repository.User, cacheDriver cache.Driver) User {
	return &user{
		db:    db,
		cache: cache.New[domain.User](cacheDriver),
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

func (u *user) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := u.db.GetByEmail(ctx, email)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}
