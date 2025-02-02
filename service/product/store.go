package product

import (
	"context"
	"fmt"

	"github.com/zechao158/ecomm/storage"
	"github.com/zechao158/ecomm/types"
	"gorm.io/gorm"
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
