package usecase

import (
	"context"
	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/4uvirik/ProductService/pkg/logger/sl"
	"github.com/go-playground/validator/v10"
	"log/slog"
)

type CategoryUseCase struct {
	oper     CategoryOperations
	logger   *slog.Logger
	validate *validator.Validate
}

func NewCategoryUseCase(oper CategoryOperations, logger *slog.Logger) *CategoryUseCase {
	return &CategoryUseCase{oper: oper, logger: logger, validate: validator.New()}
}

// ------------- Реализация CategoryOperations --------------

// CategoryCreate - метод создания новой категории
func (u *CategoryUseCase) CategoryCreate(ctx context.Context, c *entity.Category) error {
	if err := u.validate.Struct(c); err != nil {
		u.logger.Warn("validation failed for product create", sl.Err(err))
		return err
	}
	return u.oper.CategoryCreate(ctx, c)
}

// CategoryUpdate - метод изменения названия категории
func (u *CategoryUseCase) CategoryUpdate(ctx context.Context, c *entity.Category) error {
	if err := u.validate.Struct(c); err != nil {
		u.logger.Warn("validation failed for product update", sl.Err(err))
		return err
	}
	return u.oper.CategoryUpdate(ctx, c)
}

// CategoryDelete - метод удаления категории по id
func (u *CategoryUseCase) CategoryDelete(ctx context.Context, id int) error {
	if id <= 0 {
		return entity.ErrNotFound
	}
	return u.oper.CategoryDelete(ctx, id)
}
