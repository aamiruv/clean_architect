package user

import (
	"context"
	"database/sql"

	"github.com/amirzayi/clean_architect/internal/domain"
)

type userSQLRepo struct {
	db *sql.DB
}

func NewUserSQLRepository(db *sql.DB) *userSQLRepo {
	return &userSQLRepo{db: db}
}

func (r *userSQLRepo) Create(ctx context.Context, user domain.User) error {
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

func (r *userSQLRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	res := r.db.QueryRowContext(ctx, "SELECT * FROM user WHERE email=? LIMIT 1", email)

	var user domain.User
	err := res.Scan(&user)
	return user, err
}
