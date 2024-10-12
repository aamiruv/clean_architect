package sqldb

import (
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

func (r userSQLRepository) Create(auth domain.User) error {
	_, err := r.db.Exec(`INSERT INTO user 
(id,name,phone,email,password,status,created_at) 
VALUES(?,?,?,?,?,?,?)`,
		auth.ID, auth.Name, auth.PhoneNumber, auth.Email, auth.Password, auth.Status, auth.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
