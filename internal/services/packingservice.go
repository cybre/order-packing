package services

import (
	"context"
	"fmt"
	"math"
	"slices"

	"github.com/cybre/order-packing/internal/models"
)

var (
	// ErrNoPackSizesAvailable is returned when there are no pack sizes available
	ErrNoPackSizesAvailable = fmt.Errorf("no pack sizes available")

	// ErrOrderQuantity is returned when the order quantity is less than or equal to 0
	ErrOrderQuantity = fmt.Errorf("order quantity must be greater than 0")
)

// PackSizeProvider describes a type that can provide pack sizes
type PackSizeProvider interface {
	// GetPackSizes returns the available pack sizes
	GetPackSizes(ctx context.Context) ([]models.PackSize, error)
	// Update updates the pack sizes
	Update(ctx context.Context, packSizes []models.PackSize) error
}

// PackingService is a service that can calculate the number of packs required to fulfill an order
type PackingService struct {
	packSizeProvider PackSizeProvider
}

// NewPackingService returns a new PackingService with the specified pack size provider
func NewPackingService(packSizeProvider PackSizeProvider) *PackingService {
	return &PackingService{packSizeProvider}
}

// UpdatePackSizes updates the available pack sizes
func (s PackingService) UpdatePackSizes(ctx context.Context, packSizes []models.PackSize) error {
	return s.packSizeProvider.Update(ctx, packSizes)
}

// GetPackSizes returns the available pack sizes
func (s PackingService) GetPackSizes(ctx context.Context) ([]models.PackSize, error) {
	return s.packSizeProvider.GetPackSizes(ctx)
}

// CalculatePacks returns the number of packs required to fulfill the specified order
func (s PackingService) CalculatePacks(ctx context.Context, order models.Order) (map[int]int, error) {
	if order.ItemQty <= 0 {
		return nil, ErrOrderQuantity
	}

	// Get the available pack sizes
	packSizes, err := s.packSizeProvider.GetPackSizes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pack sizes: %w", err)
	}

	if len(packSizes) == 0 {
		return nil, ErrNoPackSizesAvailable
	}

	return minPacks(packSizes, order.ItemQty), nil
}

func minPacks(packSizes []models.PackSize, orderQty int) map[int]int {
	// Ensure the pack sizes are sorted in ascending order
	slices.SortFunc(packSizes, func(a, b models.PackSize) int {
		return a.MaxItems - b.MaxItems
	})

	// Calculate the maximum amount to consider, including possible overshoots (orderQty + smallest pack size)
	maxAmount := orderQty + packSizes[0].MaxItems

	// Initialize the dynamic programming table (for memoization)
	dp := make([]int, maxAmount+1)
	for i := range dp {
		dp[i] = math.MaxInt32 // Initialize with a large value
	}

	// Keep track of the pack sizes used for each amount (for backtracking)
	packSizesUsed := make([]int, maxAmount+1)

	// Base case: 0 packs are needed for 0 items
	dp[0] = 0

	// Calculate the minimum number of packs needed for each amount from packSize.MaxItems to maxAmount
	for _, packSize := range packSizes {
		for i := packSize.MaxItems; i <= maxAmount; i++ {
			if dp[i-packSize.MaxItems]+1 < dp[i] {
				dp[i] = dp[i-packSize.MaxItems] + 1
				packSizesUsed[i] = packSize.MaxItems
			}
		}
	}

	// Check for an exact match first
	if dp[orderQty] != math.MaxInt32 {
		bestAmount := orderQty

		// Backtrack to find the pack size combination for the exact amount
		packSizeCombination := make(map[int]int)
		for i := bestAmount; i > 0; i -= packSizesUsed[i] {
			packSizeCombination[packSizesUsed[i]]++
		}

		return packSizeCombination
	}

	// If no exact match, find the combination that minimizes the number of items sent out
	minItemsSent := math.MaxInt32
	bestAmount := 0
	for i := orderQty; i <= maxAmount; i++ {
		if dp[i] < math.MaxInt32 && i < minItemsSent {
			minItemsSent = i
			bestAmount = i
		}
	}

	// Backtrack to find the pack size combination for the best amount
	packSizeCombination := make(map[int]int)
	for i := bestAmount; i > 0; i -= packSizesUsed[i] {
		packSizeCombination[packSizesUsed[i]]++
	}

	return packSizeCombination
}
