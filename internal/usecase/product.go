package usecase

import (
	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/4uvirik/ProductService/pkg/logger/sl"
	"github.com/go-playground/validator/v10"
	"golang.org/x/net/context"
	"log/slog"
)

type ProductUseCase struct {
	oper     ProductOperations
	logger   *slog.Logger
	validate *validator.Validate
}

func NewProductUseCase(oper ProductOperations, logger *slog.Logger) *ProductUseCase {
	return &ProductUseCase{oper: oper, logger: logger, validate: validator.New()}
}

// ------------- Реализация ProductOperations --------------

// ProductCreate - метод создания нового продукта
func (u *ProductUseCase) ProductCreate(ctx context.Context, dto entity.ProductCreate) (*entity.Product, error) {
	if err := u.validate.Struct(dto); err != nil {
		u.logger.Warn("validation failed for product create", sl.Err(err))
		return nil, err
	}

	product := entity.Product{
		Name:       dto.Name,
		Price:      dto.Price,
		CategoryID: dto.CategoryID,
	}
	return u.oper.Create(ctx, product)
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
func (u *ProductUseCase) ProductUpdate(ctx context.Context, dto *entity.ProductUpdate) error {
	if err := u.validate.Struct(dto); err != nil {
		u.logger.Warn("validation failed for product update", sl.Err(err))
		return err
	}

	if dto.Name == nil && dto.Price == nil && dto.CategoryID == nil {
		return entity.ErrNoFields
	}

	existing, err := u.oper.GetByID(ctx, int(dto.ID))
	if err != nil {
		u.logger.Error("failed to get product", sl.Err(err))
		return err
	}

	if dto.Name != nil {
		existing.Name = *dto.Name
		u.logger.Info("update name, new name:", existing.Name)
	}
	if dto.Price != nil {
		existing.Price = *dto.Price
		u.logger.Info("update price, new price:", existing.Price)
	}
	if dto.CategoryID != nil {
		existing.CategoryID = dto.CategoryID
		u.logger.Info("update category id, new category id:", existing.CategoryID)
	}
	return u.oper.Update(ctx, existing)
}

// ProductDelete - метод удаления продукта по id
func (u *ProductUseCase) ProductDelete(ctx context.Context, id int) error {
	if id <= 0 {
		return entity.ErrNotFound
	}
	return u.oper.Delete(ctx, id)
}
