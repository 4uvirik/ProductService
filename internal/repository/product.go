package repository

import (
	"context"
	"log/slog"

	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewProductRepository(db *pgxpool.Pool, logger *slog.Logger) *ProductRepository {
	return &ProductRepository{db: db, logger: logger}
}

// ProductCreate - создает новый товар в базе данных. При успешном добавлении присваивает ID созданному объекту.
func (r *ProductRepository) ProductCreate(ctx context.Context, prod *entity.Product) (*entity.Product, error) {
	q := entity.QueryProductCreate

	err := r.db.
		QueryRow(ctx, q, prod.Name, prod.Price, prod.CategoryID).
		Scan(&prod.ID)
	if err != nil {
		r.logger.Error("failed to create product", slog.Any("err", err))
	}

	return prod, nil
}

// ProductGetAll - возвращает список всех товаров из базы данных.
func (r *ProductRepository) ProductGetAll(ctx context.Context) ([]entity.Product, error) {
	q := entity.QueryProductGetAll

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		r.logger.Error("failed to get all products", slog.Any("err", err))

		return nil, err
	}

	defer rows.Close()

	var allProducts []entity.Product

	for rows.Next() {
		var prod entity.Product
		if err := rows.Scan(&prod.ID, &prod.Name, &prod.Price, &prod.CategoryID); err != nil {
			r.logger.Error("failed to scan product row", slog.Any("err", err))

			return nil, err
		}

		allProducts = append(allProducts, prod)
	}

	if len(allProducts) == 0 {
		return nil, entity.ErrNotFound
	}

	return allProducts, nil
}

// ProductGetByID - возвращает один товар из базы данных по его ID.
func (r *ProductRepository) ProductGetByID(ctx context.Context, id int) (*entity.Product, error) {
	q := entity.QueryProductGetByID

	var prod entity.Product

	err := r.db.QueryRow(ctx, q, id).
		Scan(&prod.ID, &prod.Name, &prod.Price, &prod.CategoryID)
	if err != nil {
		r.logger.Error("failed to get product by id", slog.Int("id", id), slog.Any("err", err))

		return nil, entity.ErrNotFound
	}

	return &prod, nil
}

// ProductUpdate - обновляет поля (имя, цена, категория) существующего товара в базе данных по id.
func (r *ProductRepository) ProductUpdate(ctx context.Context, prod *entity.Product) error {
	q := entity.QueryProductUpdate

	execResult, err := r.db.Exec(ctx, q, prod.Name, prod.Price, prod.CategoryID, prod.ID)
	if err != nil {
		r.logger.Error("failed to update product", slog.Int("id", prod.ID), slog.Any("err", err))

		return err
	}

	if execResult.RowsAffected() == 0 {
		return entity.ErrNotFound
	}

	return nil
}

// ProductDelete - удаляет товар из базы данных по id.
func (r *ProductRepository) ProductDelete(ctx context.Context, id int) error {
	q := entity.QueryProductDelete

	execResult, err := r.db.Exec(ctx, q, id)
	if err != nil {
		r.logger.Error("failed to delete product", slog.Int("id", id), slog.Any("err", err))

		return err
	}

	if execResult.RowsAffected() == 0 {
		return entity.ErrNotFound
	}

	return nil
}
