// Package user holds all bootstrapper that connects repositories into a business flow related to manage users.
package user

import (
	"database/sql"

	"github.com/AmirMirzayi/clean_architecture/internal/user/adapter/repository"
	"github.com/AmirMirzayi/clean_architecture/internal/user/adapter/repository/sqldb"
	"github.com/AmirMirzayi/clean_architecture/internal/user/service"
)

func NewService(repo repository.Repository) service.UserService {
	return service.UserService{
		Repository: repo,
	}
}

func NewSQLRepository(db *sql.DB) repository.Repository {
	return sqldb.NewUserRepository(db)
}
