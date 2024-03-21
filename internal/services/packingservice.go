package services

import (
	"context"
	"fmt"
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

	// Ensure the pack sizes are sorted in descending order
	slices.SortFunc(packSizes, func(a, b models.PackSize) int {
		return b.MaxItems - a.MaxItems
	})

	// Generate all possible solutions
	solutions := generateSolutions(packSizes, order.ItemQty)

	// Pick the best solution and return it
	solution := pickBestSolution(solutions)

	return solution, nil
}

func generateSolutions(packSizes []models.PackSize, orderQty int) []map[int]int {
	solutions := []map[int]int{}

	// Iteratively and recursively generate all possible solutions
	for i := 0; i < len(packSizes); i++ {
		solutions = append(solutions, generateSolution(packSizes[i:], orderQty))
	}

	return solutions
}

func generateSolution(packSizes []models.PackSize, orderQty int) map[int]int {
	solution := map[int]int{}

	// Base case
	if len(packSizes) == 0 {
		return solution
	}

	for i := 0; i < len(packSizes); i++ {
		// Determine the number of packs of the current pack size that can be used to fulfill the order
		packSize := packSizes[i]

		packQty := orderQty / packSize.MaxItems

		// If the pack size can be used to fulfill the order, use it and generate the remainder of the solution
		// using the remaining pack sizes
		if packQty > 0 {
			solution[packSize.MaxItems] = packQty
			orderQty -= packQty * packSize.MaxItems

			branchSolution := generateSolution(packSizes[i+1:], orderQty)
			for k, v := range branchSolution {
				solution[k] += v
				orderQty -= k * v
			}

			break
		}

		// Otherwise, if this is the last pack size that can fully fulfill the order, use it and break
		if i != len(packSizes)-1 && packSizes[i+1].MaxItems < orderQty {
			solution[packSize.MaxItems] = 1
			orderQty -= packSize.MaxItems
			break
		}
	}

	// If there is a remainder, it must be fulfilled with an additional smallest pack size
	if orderQty > 0 {
		minPackSize := packSizes[len(packSizes)-1].MaxItems
		solution[minPackSize]++
	}

	return solution
}

func pickBestSolution(solutions []map[int]int) map[int]int {
	bestSolution := solutions[0]

	for _, solution := range solutions[1:] {
		// Pick the solution with the least number of items
		if totalItems(solution) < totalItems(bestSolution) {
			bestSolution = solution
		}

		// If the number of items is the same, pick the solution with the least number of packs
		if totalItems(solution) == totalItems(bestSolution) && totalPacks(solution) < totalPacks(bestSolution) {
			bestSolution = solution
		}
	}

	return bestSolution
}

func totalItems(packs map[int]int) int {
	total := 0
	for k, v := range packs {
		total += k * v
	}

	return total
}

func totalPacks(packs map[int]int) int {
	total := 0
	for _, v := range packs {
		total += v
	}

	return total
}
