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

	"github.com/cybre/re-partners-assignment/internal/models"
	"github.com/cybre/re-partners-assignment/internal/ui/templates"
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

	e.GET("/", func(c echo.Context) error {
		packSizes, err := getPackSizes(c.Request().Context(), apiAddress)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		pageData := map[string]interface{}{
			"PackSizes": packSizes,
		}

		return c.Render(http.StatusOK, "index", pageData)
	})

	e.POST("/", func(c echo.Context) error {
		orderQtyFormData := c.FormValue("quantity")
		if orderQtyFormData == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "quantity is required"})
		}

		orderQty, err := strconv.Atoi(orderQtyFormData)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		orderData, err := json.Marshal(models.Order{ItemQty: orderQty})
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

		results := []map[string]interface{}{}
		for packSize, packQty := range orderPacks {
			results = append(results, map[string]interface{}{
				"Size":     packSize,
				"Quantity": packQty,
			})
		}

		slices.SortFunc(results, func(a, b map[string]interface{}) int {
			return b["Size"].(int) - a["Size"].(int)
		})

		packSizes, err := getPackSizes(c.Request().Context(), apiAddress)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		pageData := map[string]interface{}{
			"PackSizes": packSizes,
			"Results":   results,
		}

		c.Set("results", results)

		return c.Render(http.StatusOK, "index", pageData)
	})

	e.POST("/pack-sizes", func(c echo.Context) error {
		packSizesFormData := c.FormValue("packSizes")
		if packSizesFormData == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "packSizes is required"})
		}

		packSizesStringArr := strings.Split(packSizesFormData, ",")

		packSizeModels := []models.PackSize{}
		for _, packSize := range packSizesStringArr {
			packSizeInt, err := strconv.Atoi(strings.TrimSpace(packSize))
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			}

			packSizeModels = append(packSizeModels, models.PackSize{MaxItems: packSizeInt})
		}

		packSizeData, err := json.Marshal(packSizeModels)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		resp, err := http.DefaultClient.Post(apiAddress+"/pack-sizes", "application/json", bytes.NewReader(packSizeData))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update pack sizes"})
		}

		return c.Redirect(http.StatusFound, "/")
	})

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

	packSizesView := []int{}
	for _, packSizeData := range packSizesData {
		packSizesView = append(packSizesView, packSizeData.MaxItems)
	}

	sort.Ints(packSizesView)

	return packSizesView, nil
}
