package repository

import (
	"context"
	"github.com/4uvirik/ProductService/internal/entity"
)

// ProductOperations - операции хранилища для Product
type ProductOperations interface {
	Create(ctx context.Context, p *entity.Product) error
	GetAll(ctx context.Context) ([]entity.Product, error)
	GetByID(ctx context.Context, id int) (*entity.Product, error)
	Update(ctx context.Context, p *entity.Product) error
	Delete(ctx context.Context, id int) error
}

// CategoryOperations - операции хранилища для Category
type CategoryOperations interface {
	Create(ctx context.Context, p *entity.Category) error
	Update(ctx context.Context, p *entity.Category) error
	Delete(ctx context.Context, id int) error
}
