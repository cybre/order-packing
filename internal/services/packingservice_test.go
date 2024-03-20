package services_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/cybre/re-partners-assignment/internal/models"
	"github.com/cybre/re-partners-assignment/internal/services"
	"github.com/cybre/re-partners-assignment/internal/services/testdata"
)

func TestCalculatePacls_ProviderError_ReturnError(t *testing.T) {
	t.Parallel()

	// Arrange
	expectedErr := errors.New("provider error")
	service := services.NewPackingService(&testdata.MockPackSizeProvider{Error: expectedErr})
	order := models.Order{ItemQty: 10}

	// Act
	_, err := service.CalculatePacks(context.Background(), order)

	// Assert
	if err == nil {
		t.Error("expected an error, but got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error to be %v, but got %v", expectedErr, err)
	}
}

func TestCalculatePacks(t *testing.T) {
	t.Parallel()

	defaultPackSizes := []models.PackSize{
		{MaxItems: 250},
		{MaxItems: 500},
		{MaxItems: 1000},
		{MaxItems: 2000},
		{MaxItems: 5000},
	}

	testCases := []struct {
		name          string
		packSizes     []models.PackSize
		order         models.Order
		expectedPacks map[int]int
		expectedErr   error
	}{
		{
			name:      "Order quantity is 1",
			packSizes: defaultPackSizes,
			order:     models.Order{ItemQty: 1},
			expectedPacks: map[int]int{
				250: 1,
			},
			expectedErr: nil,
		},
		{
			name:      "Order quantity is 250",
			packSizes: defaultPackSizes,
			order:     models.Order{ItemQty: 250},
			expectedPacks: map[int]int{
				250: 1,
			},
			expectedErr: nil,
		},
		{
			name:      "Order quantity is 251",
			packSizes: defaultPackSizes,
			order:     models.Order{ItemQty: 251},
			expectedPacks: map[int]int{
				500: 1,
			},
			expectedErr: nil,
		},
		{
			name:      "Order quantity is 501",
			packSizes: defaultPackSizes,
			order:     models.Order{ItemQty: 501},
			expectedPacks: map[int]int{
				500: 1,
				250: 1,
			},
			expectedErr: nil,
		},
		{
			name:      "Order quantity is 12001",
			packSizes: defaultPackSizes,
			order:     models.Order{ItemQty: 12001},
			expectedPacks: map[int]int{
				5000: 2,
				2000: 1,
				250:  1,
			},
			expectedErr: nil,
		},
		{
			name:      "Order quantity is 15001",
			packSizes: defaultPackSizes,
			order:     models.Order{ItemQty: 15001},
			expectedPacks: map[int]int{
				5000: 3,
				250:  1,
			},
			expectedErr: nil,
		},

		{
			name:      "Order quantity is 1251",
			packSizes: defaultPackSizes,
			order:     models.Order{ItemQty: 1251},
			expectedPacks: map[int]int{
				1000: 1,
				500:  1,
			},
			expectedErr: nil,
		},
		{
			name:          "Order quantity is 0",
			packSizes:     defaultPackSizes,
			order:         models.Order{ItemQty: 0},
			expectedPacks: nil,
			expectedErr:   services.ErrOrderQuantity,
		},
		{
			name:          "No pack sizes available",
			packSizes:     []models.PackSize{},
			order:         models.Order{ItemQty: 10},
			expectedPacks: nil,
			expectedErr:   services.ErrNoPackSizesAvailable,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := services.NewPackingService(&testdata.MockPackSizeProvider{
				PackSizes: tc.packSizes,
			})

			deadline, _ := t.Deadline()
			ctx, cancel := context.WithDeadline(context.Background(), deadline)
			defer cancel()

			packs, err := service.CalculatePacks(ctx, tc.order)

			if !reflect.DeepEqual(packs, tc.expectedPacks) {
				t.Errorf("expected packs to be %v, but got %v", tc.expectedPacks, packs)
			}

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error to be %v, but got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUpdatePackSizes_ProviderError_ReturnError(t *testing.T) {
	t.Parallel()

	// Arrange
	expectedErr := errors.New("provider error")
	service := services.NewPackingService(&testdata.MockPackSizeProvider{Error: expectedErr})
	packSizes := []models.PackSize{{MaxItems: 250}}

	// Act
	err := service.UpdatePackSizes(context.Background(), packSizes)

	// Assert
	if err == nil {
		t.Error("expected an error, but got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error to be %v, but got %v", expectedErr, err)
	}
}

func TestGetPackSizes_ProviderError_ReturnError(t *testing.T) {
	t.Parallel()

	// Arrange
	expectedErr := errors.New("provider error")
	service := services.NewPackingService(&testdata.MockPackSizeProvider{Error: expectedErr})

	// Act
	_, err := service.GetPackSizes(context.Background())

	// Assert
	if err == nil {
		t.Error("expected an error, but got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error to be %v, but got %v", expectedErr, err)
	}
}

func TestGetPackSizes(t *testing.T) {
	t.Parallel()

	// Arrange
	expectedPackSizes := []models.PackSize{
		{MaxItems: 250},
		{MaxItems: 500},
		{MaxItems: 1000},
		{MaxItems: 2000},
		{MaxItems: 5000},
	}
	service := services.NewPackingService(&testdata.MockPackSizeProvider{PackSizes: expectedPackSizes})

	// Act
	packSizes, err := service.GetPackSizes(context.Background())

	// Assert
	if err != nil {
		t.Fatalf("failed to get pack sizes: %v", err)
	}

	if !reflect.DeepEqual(packSizes, expectedPackSizes) {
		t.Errorf("expected pack sizes to be %v, but got %v", expectedPackSizes, packSizes)
	}
}
