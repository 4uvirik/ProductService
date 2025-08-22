package repository

import (
	"context"
	"fmt"
	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/4uvirik/ProductService/pkg"
	"github.com/4uvirik/ProductService/service/logger/sl"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Postgres struct {
	Pool   *pgxpool.Pool
	Logger *slog.Logger
}

func NewPostgres(ctx context.Context, dsn string, logger *slog.Logger) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Error("failed to connect to DB", sl.Err(err))
		return nil, fmt.Errorf("pgx connect failed: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		logger.Error("bad DB ping", sl.Err(err))
		return nil, fmt.Errorf("DB ping: %w", err)
	}

	return &Postgres{Pool: pool, Logger: logger}, nil
}

type ProductRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewProductRepository(db *pgxpool.Pool, logger *slog.Logger) *ProductRepository {
	return &ProductRepository{db: db, logger: logger}
}

// ------------- Реализация ProductOperations --------------

func (r *ProductRepository) ProductCreate(ctx context.Context, p *entity.Product) error {
	q := entity.QueryProductCreate
	err := r.db.
		QueryRow(ctx, q, p.Name, p.Price, p.CategoryID).
		Scan(&p.ID)
	if err != nil {
		r.logger.Error("failed to create product", slog.Any("err", err))
	}
	return nil
}

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
		return nil, pkg.ErrNoFound
	}
	return allProducts, nil
}

func (r *ProductRepository) ProductGetByID(ctx context.Context, id int) (*entity.Product, error) {
	q := entity.QueryProductGetByID
	var p entity.Product
	err := r.db.QueryRow(ctx, q, id).
		Scan(&p.ID, &p.Name, &p.Price, p.CategoryID)
	if err != nil {
		r.logger.Error("failed to get product by id", slog.Int("id", id), slog.Any("err", err))
		return nil, pkg.ErrNoFound
	}
	return &p, nil
}

func (r *ProductRepository) ProductUpdate(ctx context.Context, p *entity.Product) error {
	q := entity.QueryProductUpdate
	ct, err := r.db.Exec(ctx, q, p.Name, p.Price, p.CategoryID, p.ID)
	if err != nil {
		r.logger.Error("failed to update product", slog.Int("id", p.ID), slog.Any("err", err))
		return err
	}
	if ct.RowsAffected() == 0 {
		return pkg.ErrNoFound
	}
	return nil
}

func (r *ProductRepository) ProductDelete(ctx context.Context, id int) error {
	q := entity.QueryProductDelete
	ct, err := r.db.Exec(ctx, q, id)
	if err != nil {
		r.logger.Error("failed to delete product", slog.Int("id", id), slog.Any("err", err))
		return err
	}
	if ct.RowsAffected() == 0 {
		return pkg.ErrNoFound
	}
	return nil
}

// ------------- Реализация CategoryOperations --------------

func (r *ProductRepository) CategoryCreate(ctx context.Context, c *entity.Category) error {
	q := entity.QueryCategoryCreate
	err := r.db.
		QueryRow(ctx, q, c.Name).
		Scan(&c.ID)
	if err != nil {
		r.logger.Error("failed to create category", slog.Any("err", err))
	}
	return nil
}

func (r *ProductRepository) CategoryUpdate(ctx context.Context, c *entity.Category) error {
	q := entity.QueryCategoryUpdate
	ct, err := r.db.Exec(ctx, q, c.Name, c.ID)
	if err != nil {
		r.logger.Error("failed to update category", slog.Int("id", c.ID), slog.Any("err", err))
		return err
	}
	if ct.RowsAffected() == 0 {
		return pkg.ErrNoFound
	}
	return nil
}

func (r *ProductRepository) CategoryDelete(ctx context.Context, id int) error {
	q := entity.QueryCategoryDelete
	ct, err := r.db.Exec(ctx, q, id)
	if err != nil {
		r.logger.Error("failed to delete category", slog.Int("id", id), slog.Any("err", err))
		return err
	}
	if ct.RowsAffected() == 0 {
		return pkg.ErrNoFound
	}
	return nil
}
