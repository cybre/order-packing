package ui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cybre/order-packing/internal/models"
	"github.com/cybre/order-packing/internal/ui/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// StartServer starts an HTTP server on the specified address and blocks until the context is canceled.
func StartServer(ctx context.Context, address string) error {
	e := echo.New()

	e.Renderer = templates.New()

	e.Static("/static", os.Getenv("UI_STATIC_DIR"))
	e.Use(middleware.Logger())

	buildRoutes(e)

	go func() {
		if err := e.Start(address); err != nil {
			if err == http.ErrServerClosed {
				return
			}

			panic(err)
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}

func buildRoutes(e *echo.Echo) {
	apiAddress := os.Getenv("API_REMOTE_ADDRESS")

	e.GET("/", indexHandler(apiAddress))
	e.POST("/", packOrderHandler(apiAddress))
	e.POST("/pack-sizes", updatePackSizesHandler(apiAddress))

}

func getPackSizes(ctx context.Context, address string) ([]int, error) {
	packSizes, err := http.DefaultClient.Get(address + "/pack-sizes")
	if err != nil {
		return nil, fmt.Errorf("failed to get pack sizes: %w", err)
	}
	defer packSizes.Body.Close()

	if packSizes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get pack sizes: %s", packSizes.Status)
	}

	var packSizesData []models.PackSize
	if err := json.NewDecoder(packSizes.Body).Decode(&packSizesData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pack sizes: %w", err)
	}

	return mapPackSizedToViewModel(packSizesData), nil
}

func mapPackSizedToViewModel(packSizes []models.PackSize) []int {
	packSizesView := []int{}
	for _, packSizeData := range packSizes {
		packSizesView = append(packSizesView, packSizeData.MaxItems)
	}

	sort.Ints(packSizesView)

	return packSizesView
}

func indexHandler(apiAddress string) func(c echo.Context) error {
	return func(c echo.Context) error {
		packSizes, err := getPackSizes(c.Request().Context(), apiAddress)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		pageData := map[string]interface{}{
			"PackSizes": packSizes,
		}

		return c.Render(http.StatusOK, "index", pageData)
	}
}

func packOrderHandler(apiAddress string) func(c echo.Context) error {
	return func(c echo.Context) error {
		var order models.Order
		if err := c.Bind(&order); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		orderData, err := json.Marshal(order)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		resp, err := http.DefaultClient.Post(apiAddress+"/pack-order", "application/json", bytes.NewReader(orderData))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		defer resp.Body.Close()

		var orderPacks map[int]int
		if err := json.NewDecoder(resp.Body).Decode(&orderPacks); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		packSizes, err := getPackSizes(c.Request().Context(), apiAddress)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		pageData := map[string]interface{}{
			"PackSizes": packSizes,
			"Results":   mapOrderPacksToViewModel(orderPacks),
			"ItemQty":   order.ItemQty,
		}

		return c.Render(http.StatusOK, "index", pageData)
	}
}

func mapOrderPacksToViewModel(orderPacks map[int]int) []map[string]int {
	results := []map[string]int{}
	for packSize, packQty := range orderPacks {
		results = append(results, map[string]int{
			"Size":     packSize,
			"Quantity": packQty,
		})
	}

	slices.SortFunc(results, func(a, b map[string]int) int {
		return b["Size"] - a["Size"]
	})

	return results
}

func updatePackSizesHandler(apiAddress string) func(c echo.Context) error {
	return func(c echo.Context) error {
		packSizes, err := extractPackSizes(c.FormValue("packSizes"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		packSizeData, err := json.Marshal(packSizes)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		req, err := http.NewRequest(http.MethodPut, apiAddress+"/pack-sizes", bytes.NewReader(packSizeData))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update pack sizes"})
		}

		return c.Redirect(http.StatusFound, "/")
	}
}

func extractPackSizes(packSizes string) ([]models.PackSize, error) {
	if packSizes == "" {
		return nil, fmt.Errorf("pack sizes cannot be empty")
	}

	packSizeModels := []models.PackSize{}
	for _, packSize := range strings.Split(packSizes, ",") {
		packSizeInt, err := strconv.Atoi(strings.TrimSpace(packSize))
		if err != nil {
			return nil, fmt.Errorf("failed to parse pack size: %w", err)
		}

		packSizeModels = append(packSizeModels, models.PackSize{MaxItems: packSizeInt})
	}

	return packSizeModels, nil
}
