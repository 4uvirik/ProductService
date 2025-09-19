APP_NAME := Product service
ENTRYPOINT := cmd/main.go

.PHONY: tests cover integration migrate-up migrate-down migrate-status migrate-create lint lint-fix

# ---------- Unit Тесты ----------
tests:
	go test ./... -short -coverprofile=covarage.out
	go tool cover -func=covarage.out | grep total

cover:
	go tool cover -html=covarage.out -o coverage.html
	firefox coverage.html

# ---------- Integration Тесты ----------
integration:
	go test -tags=integration ./tests/integration/... -coverprofile=integration.out

# ---------- Миграции Goose ----------
DB_URL := "user=postgres password=postgres dbname=postgres host=localhost port=5432 sslmode=disable"

migrate-up:
	goose -dir ./migrations postgres $(DB_URL) up

migrate-down:
	goose -dir ./migrations postgres $(DB_URL) down

migrate-status:
	goose -dir ./migrations postgres $(DB_URL) status

migrate-create:
	goose -dir ./migrations create $(name) sql

# ---------- Проверка линтерами ----------
lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...