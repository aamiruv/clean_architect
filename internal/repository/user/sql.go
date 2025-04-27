package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/google/uuid"
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
	(id,name,phone,email,password,status,role,created_at)
	VALUES(?,?,?,?,?,?,?,?)`,
		user.ID, user.Name, user.PhoneNumber, user.Email, user.Password, user.Status, user.Role, user.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return err
	}
	return nil
}

func (r *userSQLRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM user WHERE id=? LIMIT 1", id)
	return r.scanRow(row)
}

func (r *userSQLRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM user WHERE email=? LIMIT 1", email)
	return r.scanRow(row)
}

func (r *userSQLRepo) scanRow(row *sql.Row) (domain.User, error) {
	var user domain.User
	var t string
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Email, &user.Password, &user.Status, &user.Role, &t)
	user.CreatedAt, _ = time.Parse(time.RFC3339, t)
	return user, err
}

func (r *userSQLRepo) scanRows(rows *sql.Rows) ([]domain.User, error) {
	var (
		users []domain.User
		user  domain.User
		t     string
	)
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Email, &user.Password, &user.Status, &user.Role, &t); err != nil {
			return users, err
		}
		createdAt, _ := time.Parse(time.RFC3339, t)
		user.CreatedAt = createdAt
		users = append(users, user)
	}
	return users, nil
}

func (r *userSQLRepo) List(ctx context.Context) ([]domain.User, error) {
	res, err := r.db.QueryContext(ctx, "SELECT * FROM user")
	if err != nil {
		return nil, err
	}
	return r.scanRows(res)
}

func (r *userSQLRepo) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, "UPDATE USER SET status=? WHERE id=?", domain.UserStatusDeleted, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userSQLRepo) Update(ctx context.Context, user domain.User) error {
	res, err := r.db.ExecContext(ctx, `
	UPDATE USER 
	SET name=?, phone=?, email=?, password=?
	WHERE id=?`,
		user.Name, user.PhoneNumber, user.Email, user.Password,
		user.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}
