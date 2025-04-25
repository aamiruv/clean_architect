// Package domain represents a User.
package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus uint8

const (
	UsereStatusNew UserStatus = iota + 1
	UserStatusActive
	UserStatusBanned
	UserStatusDeleted
)

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
