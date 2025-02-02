package cart

import (
	"context"

	"gorm.io/gorm"
)

type uow struct {
	db *gorm.DB
}


type UnitOfWork interface{
	Do(ctx context.Context,)
}
