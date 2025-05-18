package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/amirzayi/clean_architect/infra/migrations/model"
	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/pkg/paginate"
	"github.com/amirzayi/clean_architect/pkg/sqlutil"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userSQLRepo struct {
	db *sqlx.DB
}

func NewUserSQLRepository(db *sqlx.DB) *userSQLRepo {
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
	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM user WHERE id=? LIMIT 1", id)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	return model.ConvertUserToDomain(user), err
}

func (r *userSQLRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM user WHERE email=? LIMIT 1", email)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	return model.ConvertUserToDomain(user), err
}

func (r *userSQLRepo) List(ctx context.Context, pagination *paginate.Pagination) ([]domain.User, error) {
	users, err := sqlutil.PaginatedList[model.User](ctx, r.db, "user", pagination, map[string]string{
		"id":         "id",
		"name":       "name",
		"phone":      "phone",
		"email":      "email",
		"status":     "status",
		"role":       "role",
		"created_at": "created_at",
	})
	return model.ConvertUsersToDomains(users), err
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
		return domain.ErrUserNotFound
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
		return domain.ErrUserNotFound
	}
	return nil
}
