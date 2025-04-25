package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/internal/repository"
	"github.com/amirzayi/clean_architect/pkg/bus"
	"github.com/amirzayi/clean_architect/pkg/cache"
)

type User interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, user domain.User) error
}
type user struct {
	db       repository.User
	cache    cache.Cache[domain.User]
	eventBus bus.EventBus[domain.User]
}

func NewUserService(db repository.User, cacheDriver cache.Driver, eventDriver bus.Driver) User {
	return &user{
		db:       db,
		cache:    cache.New[domain.User](cacheDriver),
		eventBus: bus.New[domain.User](eventDriver),
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
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (u *user) List(ctx context.Context) ([]domain.User, error) {
	return u.db.List(ctx)
}

func (u *user) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.db.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (u *user) Update(ctx context.Context, user domain.User) error {
	if err := u.db.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (u *user) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := u.db.GetByID(ctx, id)
	if err != nil {
		return user, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
