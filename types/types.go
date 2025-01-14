package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	ID        uuid.UUID `json:"id,omitempty"`
	FirstName string    `json:"firstName" validate:"required"`
	LastName  string    `json:"lastName" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password,omitempty" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,min=3,max=130"`
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

func (User) TableName() string {
	return "ecom.users"
}
