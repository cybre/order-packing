package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/cybre/order-packing/internal/api"
	"github.com/cybre/order-packing/internal/providers"
	"github.com/cybre/order-packing/internal/services"
	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	godotenv.Load()

	packSizeProvider := providers.NewJSONPackSizeProvider(os.Getenv("PACKSIZES_JSON_FILE_PATH"))
	packingService := services.NewPackingService(packSizeProvider)

	// Start the server and block until the context is canceled (e.g. by pressing Ctrl+C in the terminal)
	api.StartServer(ctx, os.Getenv("API_ADDRESS"), packingService)
}
