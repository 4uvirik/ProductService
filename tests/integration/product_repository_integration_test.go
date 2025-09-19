//go:build integration
// +build integration

package integration

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/4uvirik/ProductService/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"
	"time"
)

var pool *pgxpool.Pool
var repo *repository.ProductRepository

func TestMain(m *testing.M) {

	flag.Parse()

	if testing.Short() {
		fmt.Println("short mode: skip integration tests")
		os.Exit(0)
	}

	ctx := context.Background()

	dsn, stopContainer, err := setupPostgresContainer(ctx)
	if err != nil {
		fmt.Printf("failed to start postgres container: %v\n", err)
		os.Exit(1)
	}
	defer stopContainer()

	if err := waitForDB(ctx, dsn, 60*time.Second); err != nil {
		fmt.Printf("db not ready: %v\n", err)
		os.Exit(1)
	}

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Printf("error open sql: %v\n", err)
		os.Exit(1)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		fmt.Printf("error goose set dialect: %v\n", err)
		_ = sqlDB.Close()
		os.Exit(1)
	}

	if err := goose.Up(sqlDB, "../../migrations"); err != nil {
		fmt.Printf("error goose up: %v\n", err)
		_ = sqlDB.Close()
		os.Exit(1)
	}
	_ = sqlDB.Close()

	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Printf("error pgxpool.New: %v\n", err)
		os.Exit(1)
	}

	if err := pool.Ping(ctx); err != nil {
		fmt.Printf("error pool.Ping: %v\n", err)
		os.Exit(1)
	}

	repo = repository.NewProductRepository(pool, slog.Default())

	run := m.Run()

	pool.Close()

	os.Exit(run)
}

func cleanTables(t *testing.T) {
	if pool == nil {
		t.Fatal("pgx pool is nil, TestMain probably failed")
	}
	_, err := pool.Exec(context.Background(), "TRUNCATE TABLE products, category RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

func TestProductRepositoryMethods(t *testing.T) {
	ctx := context.Background()

	cleanTables(t)
	_, err := repo.ProductGetAll(ctx)
	require.ErrorIs(t, err, entity.ErrNotFound)

	_, err = pool.Exec(ctx, `INSERT INTO category (name) VALUES ($1)`, "test category")
	require.NoError(t, err)

	var catID int
	err = pool.QueryRow(ctx, `SELECT id FROM category WHERE name=$1`, "test category").Scan(&catID)
	require.NoError(t, err)

	// Тест ProductCreate
	product := entity.Product{
		Name:       "new product",
		Price:      55.55,
		CategoryID: &catID,
	}
	createdProduct, err := repo.ProductCreate(ctx, &product)
	require.NoError(t, err)
	require.NotNil(t, createdProduct)
	require.NotZero(t, createdProduct.ID)
	require.Equal(t, product.Name, createdProduct.Name)

	// Тест ProductGetAll
	products, err := repo.ProductGetAll(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(products), 1)

	// Тест ProductGetByID
	got, err := repo.ProductGetByID(ctx, createdProduct.ID)
	require.NoError(t, err)
	require.Equal(t, createdProduct.Name, got.Name)
	require.Equal(t, createdProduct.ID, got.ID)

	// Тест ProductUpdate
	createdProduct.Name = "update product"
	createdProduct.Price = 99.99
	err = repo.ProductUpdate(ctx, createdProduct)
	require.NoError(t, err)

	updatedProduct, err := repo.ProductGetByID(ctx, createdProduct.ID)
	require.NoError(t, err)
	require.Equal(t, createdProduct.Name, updatedProduct.Name)
	require.Equal(t, createdProduct.Price, updatedProduct.Price)

	// Тест ProductDelete
	err = repo.ProductDelete(ctx, createdProduct.ID)
	require.NoError(t, err)
	_, err = repo.ProductGetByID(ctx, createdProduct.ID)
	require.ErrorIs(t, err, entity.ErrNotFound)
}
