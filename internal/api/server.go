package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cybre/re-partners-assignment/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type PackingService interface {
	CalculatePacks(context.Context, models.Order) (map[int]int, error)
	UpdatePackSizes(context.Context, []models.PackSize) error
	GetPackSizes(context.Context) ([]models.PackSize, error)
}

// StartServer starts an HTTP server on the specified address and blocks until the context is canceled.
func StartServer(ctx context.Context, address string, packingService PackingService) error {
	e := echo.New()

	buildRoutes(e, packingService)
	e.Use(middleware.Logger())

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

func buildRoutes(e *echo.Echo, packingService PackingService) {
	e.GET("/pack-sizes", func(c echo.Context) error {
		packSizes, err := packingService.GetPackSizes(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, packSizes)
	})

	e.POST("/pack-sizes", func(c echo.Context) error {
		var packSizes []models.PackSize
		if err := c.Bind(&packSizes); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		if err := packingService.UpdatePackSizes(c.Request().Context(), packSizes); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.NoContent(http.StatusNoContent)
	})

	e.POST("/pack-order", func(c echo.Context) error {
		var order models.Order
		if err := c.Bind(&order); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		packs, err := packingService.CalculatePacks(c.Request().Context(), order)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, packs)
	})
}
