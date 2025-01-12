package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RegisterUserPlayload struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

//go:generate moq -rm -pkg mocks -out mocks/user_mock.go . UserRepository:MockUserRepository
type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
}

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primarykey"`
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
