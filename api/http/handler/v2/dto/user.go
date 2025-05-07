package dto

import (
	"time"

	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Status      string    `json:"status"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

type CreateUserRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
}

func (r CreateUserRequest) ToDomain() domain.User {
	return domain.User{
		Name:        r.Name,
		PhoneNumber: r.PhoneNumber,
		Email:       r.Email,
		Password:    r.Password,
		Role:        domain.UserRole(r.Role),
	}
}

type UpdateUserRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
}

func (r UpdateUserRequest) ToDomain() domain.User {
	return domain.User{
		Name:        r.Name,
		PhoneNumber: r.PhoneNumber,
		Email:       r.Email,
		Password:    r.Password,
		Role:        domain.UserRole(r.Role),
	}
}

func UserDomainToDTO(u domain.User) UserResponse {
	return UserResponse{
		ID:          u.ID,
		Name:        u.Name,
		PhoneNumber: u.PhoneNumber,
		Email:       u.Email,
		Status:      u.Status.String(),
		Role:        string(u.Role),
		CreatedAt:   u.CreatedAt,
	}
}
