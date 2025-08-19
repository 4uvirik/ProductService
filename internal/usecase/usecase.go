package usecase

import (
	"context"
	"fmt"
	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/4uvirik/ProductService/internal/repository"
	"time"
)

type ProductUseCase struct {
	oper    repository.ProductOperations
	timeout time.Duration
}

type CategoryUseCase struct {
	oper    repository.CategoryOperations
	timeout time.Duration
}

func NewProductUseCase(oper repository.ProductOperations, timeout time.Duration) *ProductUseCase {
	return &ProductUseCase{oper: oper, timeout: timeout}
}

func NewCategoryUseCase(oper repository.CategoryOperations, timeout time.Duration) *CategoryUseCase {
	return &CategoryUseCase{oper: oper, timeout: timeout}
}

// ---------- Products ----------

func (u *ProductUseCase) ProductCreate(ctx context.Context, p *entity.Product) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	if p.Name == "" {
		return fmt.Errorf("name required")
	}
	if p.Price <= 0 {
		return fmt.Errorf("price must be > 0")
	}
	if p.CategoryID == nil {
		return fmt.Errorf("category required")
	}
	return u.oper.Create(ctx, p)
}

func (u *ProductUseCase) ProductGetAll(ctx context.Context) ([]entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.oper.GetAll(ctx)
}

func (u *ProductUseCase) ProductGetByID(ctx context.Context, id int) (*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.oper.GetByID(ctx, id)
}

func (u *ProductUseCase) ProductUpdatePrice(ctx context.Context, id int, newPrice float64) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.oper.UpdatePrice(ctx, id, newPrice)
}

func (u *ProductUseCase) ProductUpdate(ctx context.Context, p *entity.Product) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	if p.ID == 0 {
		return fmt.Errorf("id required")
	}
	if p.Name == "" {
		return fmt.Errorf("name required")
	}
	if p.Price <= 0 {
		return fmt.Errorf("price must be > 0")
	}
	if p.CategoryID == nil {
		return fmt.Errorf("category required")
	}
	return u.oper.Update(ctx, p)
}

func (u *ProductUseCase) ProductDelete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.oper.Delete(ctx, id)
}

// ---------- Category ----------

func (u *CategoryUseCase) CategoryCreate(ctx context.Context, p *entity.Category) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	if p.Name == "" {
		return fmt.Errorf("name required")
	}
	return u.oper.Create(ctx, p)
}

func (u *CategoryUseCase) CategoryUpdate(ctx context.Context, p *entity.Category) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	if p.ID == 0 {
		return fmt.Errorf("id required")
	}
	if p.Name == "" {
		return fmt.Errorf("name required")
	}
	return u.oper.Update(ctx, p)
}

func (u *CategoryUseCase) CategoryDelete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.oper.Delete(ctx, id)
}
