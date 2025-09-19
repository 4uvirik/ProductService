//go:build integration
// +build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

func setupPostgresContainer(ctx context.Context) (dsn string, stopContainer func() error, err error) {
	user := "testuser"
	password := "testpass"
	dbName := "testdbname"
	port, err := nat.NewPort("tcp", "5432")
	if err != nil {
		return "", nil, fmt.Errorf("failed to pars port: %w", err)
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": password,
			"POSTGRES_DB":       dbName,
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(port),
			wait.ForLog("database system is ready to accept connections"),
		).WithStartupTimeout(60 * time.Second),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", nil, fmt.Errorf("error start container: %w", err)
	}

	host, err := c.Host(ctx)
	if err != nil {
		_ = c.Terminate(ctx)
		return "", nil, fmt.Errorf("error get host: %w", err)
	}

	mappedPort, err := c.MappedPort(ctx, port)
	if err != nil {
		_ = c.Terminate(ctx)
		return "", nil, fmt.Errorf("error get mapped port: %w", err)
	}

	dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, mappedPort.Port(), dbName,
	)

	stopContainer = func() error {
		return c.Terminate(ctx)
	}

	return dsn, stopContainer, nil
}

func waitForDB(ctx context.Context, dsn string, timeout time.Duration) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("error sql.Open: %w", err)
	}
	defer db.Close()

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		pingCtx, cansel := context.WithTimeout(ctx, time.Second)
		err = db.PingContext(pingCtx)
		cansel()

		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		time.Sleep(time.Second)
	}

	return fmt.Errorf("timeout for db to be ready: %w", err)
}
