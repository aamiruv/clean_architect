package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/internal/repository"
	"github.com/amirzayi/clean_architect/pkg/bus"
	"github.com/amirzayi/clean_architect/pkg/cache"
	"github.com/amirzayi/clean_architect/pkg/errs"
	"github.com/amirzayi/clean_architect/pkg/paginate"
)

type User interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	List(ctx context.Context, pagination *paginate.Pagination) ([]domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, user domain.User) error
}
type user struct {
	db       repository.User
	cache    cache.Cache[domain.User]
	eventBus bus.EventBus[domain.User]
	logger   *slog.Logger
}

func NewUserService(db repository.User, cacheDriver cache.Driver, eventDriver bus.Driver, logger *slog.Logger) User {
	return &user{
		db:       db,
		cache:    cache.New[domain.User](cacheDriver),
		eventBus: bus.New[domain.User](eventDriver),
		logger:   logger,
	}
}

func (u *user) Create(ctx context.Context, user domain.User) (domain.User, error) {
	user.ID = uuid.New()
	user.Status = domain.UsereStatusNew
	user.CreatedAt = time.Now()

	if err := u.db.Create(ctx, user); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return domain.User{}, errs.New(err, errs.CodeExisted)
		}
		u.logger.Error("failed to create user", slog.Any("error", err))
		return domain.User{}, errs.New(err, errs.CodeInternal)
	}
	return user, nil
}

func (u *user) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := u.db.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.User{}, errs.NotFound("user")
		}
		u.logger.Error("failed to get user by email", slog.Any("error", err))
		return domain.User{}, errs.New(err, errs.CodeInternal)
	}
	return user, nil
}

func (u *user) List(ctx context.Context, pagination *paginate.Pagination) ([]domain.User, error) {
	users, err := u.db.List(ctx, pagination)
	if err != nil {
		u.logger.Error("failed to list users", slog.Any("error", err))
		return nil, errs.New(err, errs.CodeInternal)
	}
	return users, nil
}

func (u *user) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.db.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return errs.NotFound("user")
		}
		u.logger.Error("failed to delete user", slog.Any("error", err))
		return errs.New(err, errs.CodeInternal)
	}
	return nil
}

func (u *user) Update(ctx context.Context, user domain.User) error {
	if err := u.db.Update(ctx, user); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return errs.NotFound("user")
		}
		u.logger.Error("failed to update user", slog.Any("error", err))
		return errs.New(err, errs.CodeInternal)
	}
	return nil
}

func (u *user) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := u.db.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return user, errs.NotFound("user")
		}
		u.logger.Error("failed to get user by id", slog.Any("error", err))
		return domain.User{}, errs.New(err, errs.CodeInternal)
	}
	return user, nil
}
