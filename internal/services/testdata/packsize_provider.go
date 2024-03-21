package testdata

import (
	"context"

	"github.com/cybre/order-packing/internal/models"
)

// MockPackSizeProvider is a mock PackSizeProvider
type MockPackSizeProvider struct {
	Error     error
	PackSizes []models.PackSize
}

// GetPackSizes returns the available pack sizes or an error
func (m MockPackSizeProvider) GetPackSizes(ctx context.Context) ([]models.PackSize, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	return m.PackSizes, nil
}

// Update returns an error if one was specified
func (m MockPackSizeProvider) Update(ctx context.Context, packSizes []models.PackSize) error {
	return m.Error
}
