// Package domain represents a User.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type UserStatus uint8

const (
	UsereStatusNew UserStatus = iota + 1
	UserStatusActive
	UserStatusBanned
	UserStatusDeleted
)

func (status UserStatus) String() string {
	switch status {
	case UsereStatusNew:
		return "new"
	case UserStatusActive:
		return "active"
	case UserStatusBanned:
		return "banned"
	case UserStatusDeleted:
		return "deleted"

	default:
		return ""
	}
}

type UserRole string

const (
	UserRoleNormal UserRole = "User"
	UserRoleAdmin  UserRole = "Admin"
)

type User struct {
	ID          uuid.UUID
	Name        string
	PhoneNumber string
	Email       string
	Password    string
	Status      UserStatus
	Role        UserRole
	CreatedAt   time.Time
}
