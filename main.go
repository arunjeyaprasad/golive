package main

import (
	"log/slog"
	"os"

	"github.com/arunjeyaprasad/golive/config"
	"github.com/arunjeyaprasad/golive/server"
)

// main function initializes the application and starts the server.
// It sets up the necessary configurations, routes, and middleware.

func main() {
	config.Init()
	// Start the Server
	if err := server.StartServer(); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(-1)
	}
}
