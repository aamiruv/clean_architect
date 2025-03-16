// Package repository holds implementation for a User's Repository.
package repository

import (
	"context"

	"github.com/amirzayi/clean_architec/internal/user/domain"
)

type Repository interface {
	Create(context.Context, domain.User) error
}
