package usecase

import (
	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/4uvirik/ProductService/internal/repository"
	"github.com/4uvirik/ProductService/pkg"
	"github.com/4uvirik/ProductService/service/logger/sl"
	"github.com/go-playground/validator/v10"
	"golang.org/x/net/context"
	"log/slog"
)

type ProductUseCase struct {
	oper     repository.ProductOperations
	logger   *slog.Logger
	validate *validator.Validate
}

func NewProductUseCase(oper repository.ProductOperations, logger *slog.Logger) *ProductUseCase {
	return &ProductUseCase{oper: oper, logger: logger, validate: validator.New()}
}

// ProductCreate - метод создания нового продукта
func (u *ProductUseCase) ProductCreate(ctx context.Context, p *entity.Product) error {
	if err := u.validate.Struct(p); err != nil {
		u.logger.Warn("validation failed for product create", sl.Err(err))
		return err
	}
	return u.oper.Create(ctx, p)
}

// ProductGetAll - метод вывода всех продуктов
func (u *ProductUseCase) ProductGetAll(ctx context.Context) ([]entity.Product, error) {
	return u.oper.GetAll(ctx)
}

// ProductGetByID - метод вывода продукта по id
func (u *ProductUseCase) ProductGetByID(ctx context.Context, id int) (*entity.Product, error) {
	return u.oper.GetByID(ctx, id)
}

// ProductUpdate - метод изменения одного или нескольких полей (название, цена, категория) продукта по id
func (u *ProductUseCase) ProductUpdate(ctx context.Context, upd *entity.ProductUpdate) error {
	if err := u.validate.Struct(upd); err != nil {
		u.logger.Warn("validation failed for product update", sl.Err(err))
		return err
	}

	existing, err := u.oper.GetByID(ctx, int(upd.ID))
	if err != nil {
		u.logger.Error("failed to get product", sl.Err(err))
		return err
	}

	if upd.Name != nil {
		existing.Name = *upd.Name
		u.logger.Info("update name, new name", existing.Name)
	}
	if upd.Price != nil {
		existing.Price = *upd.Price
		u.logger.Info("update price, new price", existing.Price)
	}
	if upd.CategoryID != nil {
		existing.CategoryID = upd.CategoryID
		u.logger.Info("update category id, new category id", existing.CategoryID)
	}
	return u.oper.Update(ctx, existing)
}

// ProductDelete - метод удаления продукта по id
func (u *ProductUseCase) ProductDelete(ctx context.Context, id int) error {
	if id <= 0 {
		return pkg.ErrNoFound
	}
	return u.oper.Delete(ctx, id)
}
