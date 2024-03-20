package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/cybre/re-partners-assignment/internal/ui"
	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	godotenv.Load()

	// Start the server and block until the context is canceled (e.g. by pressing Ctrl+C in the terminal)
	ui.StartServer(ctx, os.Getenv("UI_ADDRESS"))
}
