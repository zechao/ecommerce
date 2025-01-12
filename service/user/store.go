package user

import (
	"context"
	"fmt"

	"github.com/zechao158/ecomm/storage"
	"github.com/zechao158/ecomm/types"
	"gorm.io/gorm"
)

type repository struct {
	storage.CRUDStorer[types.User]
}

func NewRepository(db *gorm.DB) types.UserRepository {
	return &repository{
		storage.New[types.User](db),
	}
}

func (s *repository) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	res, err := s.CRUDStorer.GetByFields(ctx, map[string]string{
		"email": email,
	}, false)
	if err != nil {
		return nil, fmt.Errorf("error getting user %w", err)
	}
	return res, nil
}

func (s *repository) CreateUser(ctx context.Context, user *types.User) error {
	err := s.CRUDStorer.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("error getting user %w", err)
	}
	return nil
}
