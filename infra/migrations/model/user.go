package model

import (
	"time"

	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Phone     string    `db:"phone"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Status    int       `db:"status"`
	Role      string    `db:"role"`
	CreatedAt string    `db:"created_at"`
}

func ConvertUserToDomain(user User) domain.User {
	createdAt, _ := time.Parse(time.RFC3339, user.CreatedAt)
	return domain.User{
		ID:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.Phone,
		Email:       user.Email,
		Password:    user.Password,
		Status:      domain.UserStatus(user.Status),
		Role:        domain.UserRole(user.Role),
		CreatedAt:   createdAt,
	}
}

func ConvertUsersToDomains(users []User) []domain.User {
	userDomains := make([]domain.User, 0, len(users))
	for _, user := range users {
		userDomains = append(userDomains, ConvertUserToDomain(user))
	}
	return userDomains
}
