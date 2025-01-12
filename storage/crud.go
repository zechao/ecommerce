// Package storage provides a generic implementation of CRUD (Create, Read, Update, Delete)
// operations using GORM as the underlying ORM framework. It offers a flexible and type-safe
// way to perform database operations through a consistent interface.
//
// The package implements a generic CRUDStorer interface that can work with any type,
// allowing for type-safe database operations while maintaining a clean and reusable API.
// It supports features such as:
//   - Generic CRUD operations for any entity type
//   - Support for custom query modifications through SQLModifier
//   - Row-level locking support for PostgreSQL
//   - Soft delete capability when entities include a deleted_at field
//
// Example usage:
//
//	type User struct {
//	    ID   uuid.UUID
//	    Name string
//	}
//
//	db := // initialize your GORM db
//	userStore := storage.New[User](db)
//	user, err := userStore.GetByID(ctx, someID, false)
//
//	// Create a custom repository by embedding CRUDStore
//	type UserRepository struct {
//	    db *gorm.DB
//	    CRUDStore[User]
//	}
//
//	// Add custom methods specific to User
//	func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
//	    var user User
//	    if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
//	        return nil, err
//	    }
//	    return &user, nil
//	}
//
//	// Initialize the repository
//	repo := &UserRepository{
//	    db: db,
//	    CRUDStore: New[User](db),
//	}
package storage

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SQLModifier is a type alias for a function that takes a *gorm.DB instance
// and returns a modified *gorm.DB instance. It is used to apply custom
// modifications or filters to a GORM database query.
type SQLModifier func(db *gorm.DB) *gorm.DB

// CRUDStorer defines a generic interface for basic CRUD operations.
// T represents the type of the entity that the CRUD operations will be performed on.
//
// GetAll retrieves all records that match the given SQLModifier.
// GetByID retrieves a record by its ID. The forUpdate parameter is used to lock the record for update and is only supported by PostgreSQL.
// Delete removes the specified record.
// Create adds a new record and returns the created record.
// Update modifies an existing record and returns the updated record.
type CRUDStorer[T any] interface {
	GetAll(context.Context, SQLModifier) ([]T, error)
	// GetByID retrieves a record by its ID.  for Update is used to lock the record for update. and only supportted by
	// postgresql
	GetByID(ctx context.Context, id uuid.UUID, forUpdate bool) (*T, error)
	// GetByFields retrieves a record by its field name and value.  for Update is used to lock the record for update. and only supportted by
	GetByFields(ctx context.Context, fields map[string]string, forUpdate bool) (*T, error)
	Delete(context.Context, *T) error
	//It will insert the data into the table only if the record is new,It doesn't update existing records; if the record already exists, GORM will return an error.
	Create(context.Context, *T) error
	Update(context.Context, *T) error
}

type CRUDStore[T any] struct {
	db *gorm.DB
}

// Create implements CRUDStorer.
func (c CRUDStore[T]) Create(ctx context.Context, t *T) error {
	result := c.db.WithContext(ctx).Create(&t)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			return ErrDuplicateKey
		}
		return result.Error
	}
	return nil
}

// Update implements CRUDStorer, it will create new row if the primary key of given T doesn't exist
func (c CRUDStore[T]) Update(ctx context.Context, t *T) error {
	r := c.db.Save(t)

	return r.Error
}

// Delete implements CRUDStorer, If value contains primary key it is included in the conditions. If value includes a deleted_at field, then Delete performs a soft delete instead by setting deleted_at with the current time if null.
func (c CRUDStore[T]) Delete(ctx context.Context, t *T) error {
	return c.db.Delete(t).Error
}

// GetAll implements CRUDStorer, SQLModifier allow us to add extra condition to the select statment
func (c CRUDStore[T]) GetAll(ctx context.Context, m SQLModifier) ([]T, error) {
	var results []T
	if m != nil {
		c.db = m(c.db)
	}
	if r := c.db.Find(&results); r.Error != nil {
		return nil, r.Error
	}
	return results, nil
}

// GetByID implements CRUDStorer, when forUpdate is true, it perform select for update query, which will lock the row
func (c CRUDStore[T]) GetByID(ctx context.Context, id uuid.UUID, forUpdate bool) (*T, error) {
	var r *T
	if forUpdate {
		c.db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if r := c.db.First(&r, "id = ?", id.String()); r.Error != nil {
		if r.Error == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, r.Error
	}
	return r, nil
}

// GetByID implements CRUDStorer, when forUpdate is true, it perform select for update query, which will lock the row
func (c CRUDStore[T]) GetByFields(ctx context.Context, fields map[string]string, forUpdate bool) (*T, error) {
	var r *T
	if forUpdate {
		c.db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if r := c.db.Where(fields).First(&r); r.Error != nil {
		if r.Error == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, r.Error
	}
	return r, nil
}

// New creates a new instance of CRUDStorer for the given type T.
// It takes a *gorm.DB as a parameter and returns a CRUDStorer[T].
//
// Parameters:
//   - db: A pointer to a gorm.DB instance.
//
// Returns:
//   - A CRUDStorer[T] instance initialized with the provided database connection.
func New[T any](db *gorm.DB) CRUDStorer[T] {
	return CRUDStore[T]{
		db: db,
	}
}
