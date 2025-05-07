package user

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"

	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/pkg/paginate"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type userInMemoryRepo struct {
	mu    sync.RWMutex
	store map[uuid.UUID]domain.User
}

func NewUserInMemoryRepo() *userInMemoryRepo {
	return &userInMemoryRepo{
		store: make(map[uuid.UUID]domain.User),
	}
}

func (r *userInMemoryRepo) Create(ctx context.Context, user domain.User) error {
	user, err := r.GetByID(ctx, user.ID)
	if err == nil {
		return ErrUserAlreadyExists
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[user.ID] = user
	return nil
}

func (r *userInMemoryRepo) GetByID(_ context.Context, id uuid.UUID) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.store[id]
	if !ok {
		return domain.User{}, ErrUserNotFound
	}
	return user, nil
}

func (r *userInMemoryRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.store {
		if u.Email == email {
			return u, nil
		}
	}

	return domain.User{}, ErrUserNotFound
}

func (r *userInMemoryRepo) List(_ context.Context, pagination *paginate.Pagination) ([]domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	users := make([]domain.User, 0, len(r.store))
	for _, user := range r.store {
		users = append(users, user)
	}
	return users, nil
}

func (r *userInMemoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	r.mu.Lock()
	user.Status = domain.UserStatusDeleted
	r.store[id] = user
	r.mu.Unlock()
	return nil
}

func (r *userInMemoryRepo) Update(ctx context.Context, user domain.User) error {
	u, err := r.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	r.mu.Lock()
	u.Name = user.Name
	u.PhoneNumber = user.PhoneNumber
	u.Email = user.Email
	u.Password = user.Password
	r.store[user.ID] = u
	r.mu.Unlock()
	return nil
}
