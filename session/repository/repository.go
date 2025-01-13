package repository

import (
	"context"

	"github.com/zechao158/ecomm/session"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{
		DB: db,
	}
}

func (r *User) Create(ctx context.Context, user *session.User) error {
	db := session.DB(ctx, r.DB).Omit(clause.Associations).Create(&user)
	return db.Error
}

func (r *User) CreateHistory(ctx context.Context, history *session.History) error {
	db := session.DB(ctx, r.DB).Omit(clause.Associations).Create(&history)
	return db.Error
}
