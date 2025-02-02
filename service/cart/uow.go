package cart

import (
	"gorm.io/gorm"

	"github.com/zechao158/ecomm/service/cart/order"
	orderitem "github.com/zechao158/ecomm/service/cart/order_item"
	"github.com/zechao158/ecomm/service/product"
	"github.com/zechao158/ecomm/types"
)

type OrderUOWStore struct {
	orderItemRepository types.OrderItemRepository
	orderRepository     types.OrderRepository
	productRepository   types.ProductRepository
}

type unitOfWork struct {
	db *gorm.DB
}

type UnitOfWork interface {
	Do(func(OrderUOWStore) error) error
}

func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	return &unitOfWork{db: db}
}

// Do executes the given UnitOfWorkBlock iniside a DB transaction
func (s *unitOfWork) Do(fn func(OrderUOWStore) error) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		newStore := OrderUOWStore{
			orderItemRepository: orderitem.NewRepository(tx),
			orderRepository:     order.NewRepository(tx),
			productRepository:   product.NewRepository(tx),
		}
		return fn(newStore)
	})
}
