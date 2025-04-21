package service

import (
	"github.com/amirzayi/clean_architect/internal/repository"
	"github.com/amirzayi/clean_architect/pkg/auth"
	"github.com/amirzayi/clean_architect/pkg/bus"
	"github.com/amirzayi/clean_architect/pkg/cache"
	"github.com/amirzayi/clean_architect/pkg/hash"
)

type Dependencies struct {
	Repositories *repository.Repositories
	Hasher       hash.PasswordHasher
	AuthManager  auth.Manager
	Cache        cache.Driver
	Event        bus.Driver
}

type Services struct {
	Auth Auth
	User User
}

func NewServices(deps *Dependencies) *Services {
	userService := NewUserService(deps.Repositories.User, deps.Cache, deps.Event)
	return &Services{
		User: userService,
		Auth: NewAuthService(userService, deps.Hasher, deps.AuthManager),
	}
}
