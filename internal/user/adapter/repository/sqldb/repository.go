// Package sqldb is a sql database implementation of the User's Repository interface by userSQLRepository.
package sqldb

import (
	"context"
	"database/sql"

	"github.com/AmirMirzayi/clean_architecture/internal/user/adapter/repository"
	"github.com/AmirMirzayi/clean_architecture/internal/user/domain"
)

var _ repository.Repository = userSQLRepository{}

type userSQLRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) userSQLRepository {
	return userSQLRepository{db: db}
}

func (r userSQLRepository) Create(ctx context.Context, user domain.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user 
(id,name,phone,email,password,status,created_at) 
VALUES(?,?,?,?,?,?,?)`,
		user.ID, user.Name, user.PhoneNumber, user.Email, user.Password, user.Status, user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
