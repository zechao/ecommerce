package orderitem

import (
	"github.com/zechao158/ecomm/storage"
	"github.com/zechao158/ecomm/types"
	"gorm.io/gorm"
)

type repository struct {
	storage.CRUDStorer[types.OrderItem]
}

func NewRepository(db *gorm.DB) types.OrderItemRepository {
	return &repository{
		storage.New[types.OrderItem](db),
	}
}
