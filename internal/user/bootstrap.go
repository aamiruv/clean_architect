// Package user holds all bootstrapper that connects repositories into a business flow related to manage users.
package user

import (
	"database/sql"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/amirzayi/clean_architec/internal/user/adapter/repository"
	"github.com/amirzayi/clean_architec/internal/user/adapter/repository/memory"
	"github.com/amirzayi/clean_architec/internal/user/adapter/repository/mongodb"
	"github.com/amirzayi/clean_architec/internal/user/adapter/repository/sqldb"
	"github.com/amirzayi/clean_architec/internal/user/service"
)

func NewService(repo repository.Repository) service.UserService {
	return service.UserService{
		Repository: repo,
	}
}

func NewSQLRepository(db *sql.DB) repository.Repository {
	return sqldb.NewRepository(db)
}

func NewMemoryRepository() repository.Repository {
	return memory.NewRepository()
}

func NewMongodbRepository(userCollection *mongo.Collection) repository.Repository {
	return mongodb.NewRepository(userCollection)
}
