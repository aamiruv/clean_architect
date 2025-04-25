package repository

import (
	"context"
	"database/sql"

	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/internal/repository/user"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type User interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, user domain.User) error
}

type Repositories struct {
	User User
}

func NewMongoRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		User: user.NewUserMongoRepository(db),
	}
}

func NewSQLRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User: user.NewUserSQLRepository(db),
	}
}

func NewInMemoryRepositories() *Repositories {
	return &Repositories{
		User: user.NewUserInMemoryRepo(),
	}
}
