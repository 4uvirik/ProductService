package server

import (
	"context"
	"github.com/4uvirik/ProductService/config"
	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/4uvirik/ProductService/internal/usecase"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type mockCategoryRepo struct{}

func (m *mockCategoryRepo) CategoryUpdate(ctx context.Context, cat *entity.Category) error {
	return nil
}

func (m *mockCategoryRepo) CategoryCreate(ctx context.Context, cat *entity.Category) error {
	return nil
}

func (m *mockCategoryRepo) CategoryDelete(ctx context.Context, id int) error {
	return nil
}

type mockProductRepo struct{}

func (m *mockProductRepo) ProductCreate(ctx context.Context, p entity.Product) (*entity.Product, error) {
	p.ID = 1
	return &p, nil
}
func (m *mockProductRepo) ProductGetAll(ctx context.Context) ([]entity.Product, error) {
	return []entity.Product{{ID: 1, Name: "test product", Price: 100, CategoryID: nil}}, nil
}
func (m *mockProductRepo) ProductGetByID(ctx context.Context, id int) (*entity.Product, error) {
	return &entity.Product{ID: id, Name: "test product", Price: 100, CategoryID: nil}, nil
}
func (m *mockProductRepo) ProductUpdate(ctx context.Context, p *entity.Product) error {
	return nil
}
func (m *mockProductRepo) ProductDelete(ctx context.Context, id int) error {
	return nil
}

func TestRun(t *testing.T) {

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		want   int
	}{
		{"get product by id", http.MethodGet, "/product/1", "", http.StatusOK},
		{"create product", http.MethodPost, "/product", `{"name": "test product", "price": 100}`, http.StatusCreated},
		{"update category", http.MethodPut, "/category/1", `{"name": "new category"}`, http.StatusOK},
		{"create category", http.MethodPost, "/category", `{"name": "new category"}`, http.StatusCreated},
		{"unknown rout", http.MethodGet, "/unknown", "", http.StatusNotFound},
	}

	cfg := &config.Config{
		App: config.AppConfig{
			Host: "127.0.0.1",
			Port: "8080",
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	categoryUC := usecase.NewCategoryUseCase(&mockCategoryRepo{}, logger)
	productUC := usecase.NewProductUseCase(&mockProductRepo{}, logger)

	e := Run(cfg, categoryUC, productUC, logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyReader *strings.Reader
			if tt.body != "" {
				bodyReader = strings.NewReader(tt.body)
			} else {
				bodyReader = strings.NewReader("")
			}

			req := httptest.NewRequest(tt.method, tt.path, bodyReader)
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != tt.want {
				t.Errorf("%s: got status %d, want %d", tt.name, rec.Code, tt.want)
			}

		})
	}
}
