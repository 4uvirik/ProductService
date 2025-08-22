package usecase

import (
	"context"
	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/4uvirik/ProductService/internal/repository"
	"github.com/4uvirik/ProductService/pkg"
	"github.com/4uvirik/ProductService/service/logger/sl"
	"github.com/go-playground/validator/v10"
	"log/slog"
)

type CategoryUseCase struct {
	oper     repository.CategoryOperations
	logger   *slog.Logger
	validate *validator.Validate
}

func NewCategoryUseCase(oper repository.CategoryOperations, logger *slog.Logger) *CategoryUseCase {
	return &CategoryUseCase{oper: oper, logger: logger, validate: validator.New()}
}

// ---------- Category ----------

func (u *CategoryUseCase) CategoryCreate(ctx context.Context, c *entity.Category) error {
	if err := u.validate.Struct(c); err != nil {
		u.logger.Warn("validation failed for product create", sl.Err(err))
		return err
	}
	return u.oper.Create(ctx, c)
}

func (u *CategoryUseCase) CategoryUpdate(ctx context.Context, c *entity.Category) error {
	if err := u.validate.Struct(c); err != nil {
		u.logger.Warn("validation failed for product update", sl.Err(err))
		return err
	}
	return u.oper.Update(ctx, c)
}

func (u *CategoryUseCase) CategoryDelete(ctx context.Context, id int) error {
	if id <= 0 {
		return pkg.ErrNoFound
	}
	return u.oper.Delete(ctx, id)
}
