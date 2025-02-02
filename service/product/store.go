package product

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/zechao158/ecomm/storage"
	"github.com/zechao158/ecomm/types"
)

type repository struct {
	storage.CRUDStorer[types.Product]
}

func NewRepository(db *gorm.DB) types.ProductRepository {
	return &repository{
		storage.New[types.Product](db),
	}
}

func (s *repository) GetProducts(ctx context.Context) ([]types.Product, error) {
	res, err := s.GetAll(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting all products %w", err)
	}
	return res, nil
}

func (s *repository) GetProductsByIDs(ctx context.Context, ids []uuid.UUID) ([]types.Product, error) {
	res, err := s.GetAll(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("id in (?)", ids)
	})
	if err != nil {
		return nil, fmt.Errorf("error getting all products %w", err)
	}
	return res, nil
}
