package repository

import (
	"database/sql"
	"github.com/AmirMirzayi/clean_architecture/internal/auth/domain"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return AuthRepository{db: db}
}

func (r AuthRepository) Register(auth domain.Auth) error {
	res, err := r.db.Exec("INSERT INTO user(user_name,password) VALUES(?,?)", auth.UserName, auth.Password)
	_ = res
	if err != nil {
		return err
	}
	return nil
}
