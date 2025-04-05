package user

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"

	"github.com/amirzayi/clean_architect/internal/domain"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type userInMemoryRepo struct {
	mu    sync.RWMutex
	store map[uuid.UUID]domain.User
}

func NewUserInMemoryRepo() *userInMemoryRepo {
	return &userInMemoryRepo{
		mu:    sync.RWMutex{},
		store: make(map[uuid.UUID]domain.User),
	}
}

func (r *userInMemoryRepo) Create(ctx context.Context, user domain.User) error {
	_, err := r.FindByID(ctx, user.ID)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[user.ID] = user
	return nil
}

func (r *userInMemoryRepo) FindByID(_ context.Context, id uuid.UUID) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.store[id]
	if !ok {
		return domain.User{}, ErrUserAlreadyExists
	}
	return user, nil
}
