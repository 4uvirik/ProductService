package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/4uvirik/ProductService/internal/logger/sl"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

var ErrNoFound = errors.New("no found")

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

// ------------- Реализация ProductOperations --------------

func (pg *Postgres) ProductCreate(ctx context.Context, p *entity.Product) error {
	q := entity.QueryProductCreate
	return pg.Pool.
		QueryRow(ctx, q, p.Name, p.Price, p.CategoryID).
		Scan(&p.ID)
}

func (pg *Postgres) ProductGetAll(ctx context.Context) ([]entity.Product, error) {
	q := entity.QueryProductGetAll
	rows, err := pg.Pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allProducts []entity.Product
	for rows.Next() {
		var p entity.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CategoryID); err != nil {
			return nil, err
		}
		allProducts = append(allProducts, p)
	}
	return allProducts, rows.Err()
}

func (pg *Postgres) ProductGetByID(ctx context.Context, id int) (*entity.Product, error) {
	q := entity.QueryProductGetByID
	var p entity.Product
	err := pg.Pool.QueryRow(ctx, q, id).
		Scan(&p.ID, &p.Name, &p.Price, p.CategoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoFound
		}
		return nil, err
	}
	return &p, nil
}

func (pg *Postgres) ProductUpdatePrice(ctx context.Context, id int, newPrice float64) error {
	q := entity.QueryProductUpdatePrice
	ct, err := pg.Pool.Exec(ctx, q, newPrice, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNoFound
	}
	return nil
}

func (pg *Postgres) ProductUpdate(ctx context.Context, p *entity.Product) error {
	q := entity.QueryProductUpdate
	ct, err := pg.Pool.Exec(ctx, q, p.Name, p.Price, p.CategoryID, p.ID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNoFound
	}
	return nil
}

func (pg *Postgres) ProductDelete(ctx context.Context, id int) error {
	q := entity.QueryProductDelete
	ct, err := pg.Pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNoFound
	}
	return nil
}

// ------------- Реализация CategoryOperations --------------

func (pg *Postgres) CategoryCreate(ctx context.Context, p *entity.Product) error {
	q := entity.QueryCategoryCreate
	return pg.Pool.
		QueryRow(ctx, q, p.Name).
		Scan(&p.ID)
}

func (pg *Postgres) CategoryUpdate(ctx context.Context, p *entity.Product) error {
	q := entity.QueryCategoryUpdate
	ct, err := pg.Pool.Exec(ctx, q, p.Name, p.ID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNoFound
	}
	return nil
}

func (pg *Postgres) CategoryDelete(ctx context.Context, id int) error {
	q := entity.QueryCategoryDelete
	ct, err := pg.Pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNoFound
	}
	return nil
}
