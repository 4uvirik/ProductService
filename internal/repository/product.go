package repository

import (
	"context"
	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type ProductRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewProductRepository(db *pgxpool.Pool, logger *slog.Logger) *ProductRepository {
	return &ProductRepository{db: db, logger: logger}
}

// ProductCreate - создает новый товар в базе данных. При успешном добавлении присваивает ID созданному объекту.
func (r *ProductRepository) ProductCreate(ctx context.Context, p entity.Product) (*entity.Product, error) {
	q := entity.QueryProductCreate
	err := r.db.
		QueryRow(ctx, q, p.Name, p.Price, p.CategoryID).
		Scan(&p.ID)
	if err != nil {
		r.logger.Error("failed to create product", slog.Any("err", err))
	}
	return &p, nil
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
		var p entity.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CategoryID); err != nil {
			r.logger.Error("failed to scan product row", slog.Any("err", err))
			return nil, err
		}
		allProducts = append(allProducts, p)
	}

	if len(allProducts) == 0 {
		return nil, entity.ErrNotFound
	}
	return allProducts, nil
}

// ProductGetByID - возвращает один товар из базы данных по его ID.
func (r *ProductRepository) ProductGetByID(ctx context.Context, id int) (*entity.Product, error) {
	q := entity.QueryProductGetByID
	var p entity.Product
	err := r.db.QueryRow(ctx, q, id).
		Scan(&p.ID, &p.Name, &p.Price, p.CategoryID)
	if err != nil {
		r.logger.Error("failed to get product by id", slog.Int("id", id), slog.Any("err", err))
		return nil, entity.ErrNotFound
	}
	return &p, nil
}

// ProductUpdate - обновляет поля (имя, цена, категория) существующего товара в базе данных по id.
func (r *ProductRepository) ProductUpdate(ctx context.Context, p *entity.Product) error {
	q := entity.QueryProductUpdate
	ct, err := r.db.Exec(ctx, q, p.Name, p.Price, p.CategoryID, p.ID)
	if err != nil {
		r.logger.Error("failed to update product", slog.Int("id", p.ID), slog.Any("err", err))
		return err
	}
	if ct.RowsAffected() == 0 {
		return entity.ErrNotFound
	}
	return nil
}

// ProductDelete - удаляет товар из базы данных по id.
func (r *ProductRepository) ProductDelete(ctx context.Context, id int) error {
	q := entity.QueryProductDelete
	ct, err := r.db.Exec(ctx, q, id)
	if err != nil {
		r.logger.Error("failed to delete product", slog.Int("id", id), slog.Any("err", err))
		return err
	}
	if ct.RowsAffected() == 0 {
		return entity.ErrNotFound
	}
	return nil
}
