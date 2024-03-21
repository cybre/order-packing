package testdata

import (
	"context"

	"github.com/cybre/order-packing/internal/models"
)

type MockPackingService struct {
	Error     error
	PackSizes []models.PackSize
	Packs     map[int]int
}

func (m MockPackingService) GetPackSizes(ctx context.Context) ([]models.PackSize, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	return m.PackSizes, nil
}

func (m MockPackingService) CalculatePacks(context.Context, models.Order) (map[int]int, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	return m.Packs, nil
}

func (m MockPackingService) UpdatePackSizes(ctx context.Context, packSizes []models.PackSize) error {
	return m.Error
}
