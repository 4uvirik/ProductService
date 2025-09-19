package usecase

import (
	"context"
	"github.com/4uvirik/ProductService/internal/entity"
	"github.com/stretchr/testify/mock"
)

// ---------- Mock интерфейса ProductOperations для тестов ----------
type MockProductOperations struct {
	mock.Mock
}

func (m *MockProductOperations) ProductCreate(ctx context.Context, p *entity.Product) (*entity.Product, error) {
	args := m.Called(ctx, p)
	if prod, ok := args.Get(0).(*entity.Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductOperations) ProductGetAll(ctx context.Context) ([]entity.Product, error) {
	args := m.Called(ctx)
	if prod, ok := args.Get(0).([]entity.Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductOperations) ProductGetByID(ctx context.Context, id int) (*entity.Product, error) {
	args := m.Called(ctx, id)
	if prod, ok := args.Get(0).(*entity.Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductOperations) ProductUpdate(ctx context.Context, p *entity.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProductOperations) ProductDelete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---------- Mock интерфейса CategoryOperations для тестов ----------
type MockCategoryOperations struct {
	mock.Mock
}

func (m *MockCategoryOperations) CategoryCreate(ctx context.Context, p *entity.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockCategoryOperations) CategoryUpdate(ctx context.Context, p *entity.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockCategoryOperations) CategoryDelete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
