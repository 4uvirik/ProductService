package repository

import (
	"context"
	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/4uvirik/ProductService/pkg"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type CategoryRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewCategoryRepository(db *pgxpool.Pool, logger *slog.Logger) *ProductRepository {
	return &ProductRepository{db: db, logger: logger}
}

// ------------- Реализация CategoryOperations --------------

func (r *CategoryRepository) CategoryCreate(ctx context.Context, c *entity.Category) error {
	q := entity.QueryCategoryCreate
	err := r.db.
		QueryRow(ctx, q, c.Name).
		Scan(&c.ID)
	if err != nil {
		r.logger.Error("failed to create category", slog.Any("err", err))
	}
	return nil
}

func (r *CategoryRepository) CategoryUpdate(ctx context.Context, c *entity.Category) error {
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

func (r *CategoryRepository) CategoryDelete(ctx context.Context, id int) error {
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
