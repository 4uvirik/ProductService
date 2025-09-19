package usecase

import (
	"context"
	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"testing"
)

// TestProductUseCase_ProductCreate - Unit-тест метода ProductCreate
func TestProductUseCase_ProductCreate(t *testing.T) {
	catID := 1

	tests := []struct {
		name    string
		dto     entity.ProductCreate
		mock    func(m *MockProductOperations)
		wantErr bool
	}{
		{
			name: "success",
			dto:  entity.ProductCreate{Name: "new product", Price: 10, CategoryID: &catID},
			mock: func(m *MockProductOperations) {
				m.On("ProductCreate", mock.Anything, mock.AnythingOfType("*entity.Product")).
					Return(&entity.Product{Name: "new product", Price: 10, CategoryID: &catID}, nil)
			},
			wantErr: false,
		},
		{
			name:    "validate failed (empty name)",
			dto:     entity.ProductCreate{Name: "", Price: 10, CategoryID: &catID},
			mock:    func(m *MockProductOperations) {},
			wantErr: true,
		},
		{
			name: "repository error",
			dto:  entity.ProductCreate{Name: "new product", Price: 10, CategoryID: &catID},
			mock: func(m *MockProductOperations) {
				m.On("ProductCreate", mock.Anything, mock.AnythingOfType("*entity.Product")).
					Return(nil, entity.ErrDBError)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockProductOperations{}
			tt.mock(m)
			n := NewProductUseCase(m, slog.Default())

			_, err := n.ProductCreate(context.Background(), tt.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			m.AssertExpectations(t)
		})
	}
}

// TestProductUseCase_ProductGetAll - Unit-тест метода ProductGetAll
func TestProductUseCase_ProductGetAll(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(m *MockProductOperations)
		wantErr bool
	}{
		{
			name: "success",
			mock: func(m *MockProductOperations) {
				m.On("ProductGetAll", mock.Anything).
					Return([]entity.Product{{ID: 1, Name: "some product"}}, nil)
			},
			wantErr: false,
		},
		{
			name: "repo error",
			mock: func(m *MockProductOperations) {
				m.On("ProductGetAll", mock.Anything).
					Return(nil, entity.ErrDBError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockProductOperations{}
			tt.mock(m)
			u := NewProductUseCase(m, slog.Default())

			_, err := u.ProductGetAll(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			m.AssertExpectations(t)
		})
	}
}

// TestProductUseCase_ProductGetByID - Unit-тест метода ProductGetByID
func TestProductUseCase_ProductGetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		mock    func(m *MockProductOperations)
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
			mock: func(m *MockProductOperations) {
				m.On("ProductGetByID", mock.Anything, 1).
					Return(&entity.Product{ID: 1, Name: "some product"}, nil)
			},
			wantErr: false,
		},
		{
			name: "repo error",
			id:   99,
			mock: func(m *MockProductOperations) {
				m.On("ProductGetByID", mock.Anything, 99).
					Return(nil, entity.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockProductOperations{}
			tt.mock(m)
			u := NewProductUseCase(m, slog.Default())

			_, err := u.ProductGetByID(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			m.AssertExpectations(t)
		})
	}
}

// TestProductUseCase_ProductUpdate - Unit-тест метода ProductUpdate
func TestProductUseCase_ProductUpdate(t *testing.T) {
	name := "new name"
	price := 99.99
	catID := 5

	tests := []struct {
		name    string
		dto     *entity.ProductUpdate
		mock    func(m *MockProductOperations)
		wantErr bool
	}{
		{
			name:    "validation failed (ID = 0)",
			dto:     &entity.ProductUpdate{ID: 0},
			mock:    func(m *MockProductOperations) {},
			wantErr: true,
		},
		{
			name:    "all fields nil",
			dto:     &entity.ProductUpdate{ID: 1},
			mock:    func(m *MockProductOperations) {},
			wantErr: true,
		},
		{
			name: "get product error",
			dto:  &entity.ProductUpdate{ID: 2, Name: &name},
			mock: func(m *MockProductOperations) {
				m.On("ProductGetByID", mock.Anything, 2).
					Return(nil, entity.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "success update",
			dto:  &entity.ProductUpdate{ID: 1, Name: &name, Price: &price, CategoryID: &catID},
			mock: func(m *MockProductOperations) {
				m.On("ProductGetByID", mock.Anything, 1).
					Return(&entity.Product{ID: 1, Name: "old name", Price: 10, CategoryID: &catID}, nil)
				m.On("ProductUpdate", mock.Anything, mock.AnythingOfType("*entity.Product")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "update repo error",
			dto:  &entity.ProductUpdate{ID: 1, Name: &name},
			mock: func(m *MockProductOperations) {
				m.On("ProductGetByID", mock.Anything, 1).
					Return(&entity.Product{ID: 1, Name: "old name"}, nil)
				m.On("ProductUpdate", mock.Anything, mock.AnythingOfType("*entity.Product")).
					Return(entity.ErrDBError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockProductOperations{}
			tt.mock(m)
			u := NewProductUseCase(m, slog.Default())

			err := u.ProductUpdate(context.Background(), tt.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			m.AssertExpectations(t)
		})
	}
}

// TestProductUseCase_ProductDelete -  Unit-тест метода ProductDelete
func TestProductUseCase_ProductDelete(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		mock    func(m *MockProductOperations)
		wantErr bool
	}{
		{
			name:    "invalid id",
			id:      0,
			mock:    func(m *MockProductOperations) {},
			wantErr: true,
		},
		{
			name: "success delete",
			id:   1,
			mock: func(m *MockProductOperations) {
				m.On("ProductDelete", mock.Anything, 1).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repo error",
			id:   2,
			mock: func(m *MockProductOperations) {
				m.On("ProductDelete", mock.Anything, 2).
					Return(entity.ErrDBError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockProductOperations{}
			tt.mock(m)
			u := NewProductUseCase(m, slog.Default())

			err := u.ProductDelete(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			m.AssertExpectations(t)
		})
	}
}
