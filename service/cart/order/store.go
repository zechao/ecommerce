package order

import (
	"github.com/zechao158/ecomm/storage"
	"github.com/zechao158/ecomm/types"
	"gorm.io/gorm"
)

type repository struct {
	storage.CRUDStorer[types.Order]
}

func NewRepository(db *gorm.DB) types.OrderRepository {
	return &repository{
		storage.New[types.Order](db),
	}
}
