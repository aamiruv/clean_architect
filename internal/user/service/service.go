package service

import (
	"context"
	"fmt"
	"time"

	"github.com/AmirMirzayi/clean_architecture/internal/user/adapter/repository"
	"github.com/AmirMirzayi/clean_architecture/internal/user/domain"
	"github.com/google/uuid"
)

type UserService struct {
	Repository repository.Repository
}

func (s UserService) Create(ctx context.Context, user domain.User) (domain.User, error) {
	user.ID = uuid.New()
	user.Status = domain.New
	user.CreatedAt = time.Now()
	err := s.Repository.Create(user)
	if err != nil {
		return user, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}
