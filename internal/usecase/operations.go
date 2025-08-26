package usecase

import (
	"context"
	"github.com/4uvirik/ProductService/internal/entity"
)

// ProductOperations - операции хранилища для Product
type ProductOperations interface {
	ProductCreate(ctx context.Context, p entity.Product) (*entity.Product, error)
	ProductGetAll(ctx context.Context) ([]entity.Product, error)
	ProductGetByID(ctx context.Context, id int) (*entity.Product, error)
	ProductUpdate(ctx context.Context, p *entity.Product) error
	ProductDelete(ctx context.Context, id int) error
}

// CategoryOperations - операции хранилища для Category
type CategoryOperations interface {
	CategoryCreate(ctx context.Context, p *entity.Category) error
	CategoryUpdate(ctx context.Context, p *entity.Category) error
	CategoryDelete(ctx context.Context, id int) error
}
