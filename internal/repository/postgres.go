package repository

import (
	"context"
	"fmt"
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
