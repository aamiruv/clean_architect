package dto

type CreateUserRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Password    string `json:"password"`
}
