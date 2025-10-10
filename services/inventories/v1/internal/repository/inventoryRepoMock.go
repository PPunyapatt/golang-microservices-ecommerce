package repository

import (
	"context"
	"inventories/v1/internal/constant"

	"github.com/stretchr/testify/mock"
)

type InventoryRepositoryMock struct {
	mock.Mock
}

func NewInventoryRepositoryMock() *InventoryRepositoryMock {
	return &InventoryRepositoryMock{}
}

func (m *InventoryRepositoryMock) ReserveStock(ctx context.Context, inventory *constant.Inventory) error {
	args := m.Called(ctx, inventory)
	return args.Error(0)
}
