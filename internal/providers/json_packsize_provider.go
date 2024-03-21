package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cybre/order-packing/internal/models"
)

// JSONPackSizeProvider is a PackSizeProvider that reads and stores pack sizes in a JSON file
type JSONPackSizeProvider struct {
	filePath string
}

// NewJSONPackSizeProvider returns a new JSONPackSizeProvider with the specified file path
func NewJSONPackSizeProvider(filePath string) *JSONPackSizeProvider {
	return &JSONPackSizeProvider{filePath}
}

// GetPackSizes returns the available pack sizes
func (p JSONPackSizeProvider) GetPackSizes(ctx context.Context) ([]models.PackSize, error) {
	packSizes := []models.PackSize{}

	file, err := os.Open(p.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&packSizes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file: %w", err)
	}

	return packSizes, nil
}

// Update updates the pack sizes in the JSON file
func (p JSONPackSizeProvider) Update(ctx context.Context, packSizes []models.PackSize) error {
	file, err := os.Create(p.filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(packSizes); err != nil {
		return fmt.Errorf("failed to marshal pack sizes: %w", err)
	}

	return nil
}
