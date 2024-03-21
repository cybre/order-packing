package providers_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/cybre/order-packing/internal/models"
	"github.com/cybre/order-packing/internal/providers"
)

func TestJSONPackSizeProvider_GetPackSizes(t *testing.T) {
	// Create a temporary file for testing
	file, err := os.CreateTemp("", "test_pack_sizes.json")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Write test data to the temporary file
	testData := []models.PackSize{
		{MaxItems: 10},
		{MaxItems: 20},
		{MaxItems: 30},
	}
	if err := json.NewEncoder(file).Encode(testData); err != nil {
		t.Fatalf("failed to write test data to file: %v", err)
	}

	// Create an instance of JSONPackSizeProvider
	p := providers.NewJSONPackSizeProvider(file.Name())

	// Act
	packSizes, err := p.GetPackSizes(context.Background())
	if err != nil {
		t.Fatalf("failed to get pack sizes: %v", err)
	}

	// Assert
	expectedPackSizes := []models.PackSize{
		{MaxItems: 10},
		{MaxItems: 20},
		{MaxItems: 30},
	}
	if len(packSizes) != len(expectedPackSizes) {
		t.Fatalf("unexpected number of pack sizes, got %d, want %d", len(packSizes), len(expectedPackSizes))
	}
	for i, ps := range packSizes {
		if ps != expectedPackSizes[i] {
			t.Fatalf("unexpected pack size at index %d, got %+v, want %+v", i, ps, expectedPackSizes[i])
		}
	}
}

func TestJSONPackSizeProvider_Update(t *testing.T) {
	// Create a temporary file for testing
	file, err := os.CreateTemp("", "test_pack_sizes.json")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Create an instance of JSONPackSizeProvider
	p := providers.NewJSONPackSizeProvider(file.Name())

	// Act
	err = p.Update(context.Background(), []models.PackSize{
		{MaxItems: 10},
		{MaxItems: 20},
		{MaxItems: 30},
	})
	if err != nil {
		t.Fatalf("failed to update pack sizes: %v", err)
	}

	// Read the temporary file
	file, err = os.Open(file.Name())
	if err != nil {
		t.Fatalf("failed to open temporary file: %v", err)
	}
	defer file.Close()

	// Assert
	var packSizes []models.PackSize
	if err := json.NewDecoder(file).Decode(&packSizes); err != nil {
		t.Fatalf("failed to read pack sizes from file: %v", err)
	}
	expectedPackSizes := []models.PackSize{
		{MaxItems: 10},
		{MaxItems: 20},
		{MaxItems: 30},
	}
	if len(packSizes) != len(expectedPackSizes) {
		t.Fatalf("unexpected number of pack sizes, got %d, want %d", len(packSizes), len(expectedPackSizes))
	}
	for i, ps := range packSizes {
		if ps != expectedPackSizes[i] {
			t.Fatalf("unexpected pack size at index %d, got %+v, want %+v", i, ps, expectedPackSizes[i])
		}
	}
}
