package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status uint8

const (
	New Status = iota + 1
	Active
)

type User struct {
	ID          uuid.UUID
	Name        string
	PhoneNumber string
	Email       string
	Password    string
	Status      Status
	CreatedAt   time.Time
}
