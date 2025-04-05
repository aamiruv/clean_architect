package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password string) (hashed string, err error)
}

type bcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) PasswordHasher {
	return &bcryptHasher{cost: cost}
}

func (b *bcryptHasher) Hash(pwd string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(pwd), b.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(password), nil
}
