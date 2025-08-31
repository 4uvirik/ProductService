APP_NAME := Product service
ENTRYPOINT := cmd/main.go

.PHONY: tests cover migrate-up migrate-down migrate-status migrate-create

# ---------- Unit Тесты ----------
tests:
	go test ./... -coverprofile=covarage.out
	go tool cover -func=covarage.out | grep total

cover:
	go tool cover -html=covarage.out -o coverage.html
	firefox coverage.html

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
