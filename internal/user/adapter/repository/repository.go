// Package repository holds implementation for a User's Repository.
package repository

import (
	"context"

	"github.com/AmirMirzayi/clean_architecture/internal/user/domain"
)

type Repository interface {
	Create(context.Context, domain.User) error
}
