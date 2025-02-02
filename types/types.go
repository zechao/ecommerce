package types

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/zechao158/ecomm/storage"
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
	storage.CRUDStorer[User]
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primarykey"`
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func (User) TableName() string {
	return "ecom.users"
}

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primarykey"`
	Name        string
	Description string
	Image       string
	Price       float64
	Quantity    int
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

//go:generate moq -rm -pkg mocks -out mocks/product_mock.go . ProductRepository:MockProductRepository
type ProductRepository interface {
	storage.CRUDStorer[Product]
	GetProductsByIDs(context.Context, []uuid.UUID) ([]Product, error)
}

func (Product) TableName() string {
	return "ecom.products"
}

type Order struct {
	ID        uuid.UUID `gorm:"type:uuid;primarykey"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	Total     float64
	Status    string
	Address   string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func (Order) TableName() string {
	return "ecom.orders"
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primarykey"`
	OrderID   uuid.UUID `gorm:"type:uuid"`
	ProductID uuid.UUID `gorm:"type:uuid"`
	Quantity  int
	Price     float64
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func (Order) OrderItem() string {
	return "ecom.order_items"
}

//go:generate moq -rm -pkg mocks -out mocks/product_mock.go . ProductRepository:MockProductRepository
type OrderRepository interface {
	storage.CRUDStorer[Order]
}

//go:generate moq -rm -pkg mocks -out mocks/product_mock.go . ProductRepository:MockProductRepository
type OrderItemRepository interface {
	storage.CRUDStorer[OrderItem]
}

type CartItem struct {
	ProductID uuid.UUID `gorm:"type:uuid"`
	Quantity  int
}

type CartCheckoutPayload struct {
	Items []CartItem `json:"items" validate:"required"`
}
