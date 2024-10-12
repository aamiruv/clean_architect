package repository

import (
	"github.com/AmirMirzayi/clean_architecture/internal/user/domain"
)

type Repository interface {
	Create(domain.User) error
}
