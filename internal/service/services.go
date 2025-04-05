package service

import (
	"github.com/amirzayi/clean_architect/internal/repository"
	"github.com/amirzayi/clean_architect/pkg/hash"
)

type Dependencies struct {
	Repositories *repository.Repositories
	Hasher       hash.PasswordHasher
}

type Services struct {
	Auth Auth
	User User
}

func NewServices(deps *Dependencies) *Services {
	userService := NewUserService(deps.Repositories.User)
	return &Services{
		User: userService,
		Auth: NewAuthService(userService, deps.Hasher),
	}
}
