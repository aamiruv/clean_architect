package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"

	"github.com/amirzayi/clean_architec/internal/user/adapter/repository"
	"github.com/amirzayi/clean_architec/internal/user/domain"
)

var _ repository.Repository = (*userInMemoryRepository)(nil)

type userInMemoryRepository struct {
	mu    sync.Mutex
	store map[uuid.UUID]domain.User
}

func NewRepository() repository.Repository {
	return &userInMemoryRepository{
		mu:    sync.Mutex{},
		store: make(map[uuid.UUID]domain.User),
	}
}

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

func (r *userInMemoryRepository) Create(ctx context.Context, user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[user.ID]; ok {
		return ErrUserAlreadyExists
	}
	r.store[user.ID] = user
	return nil
}
