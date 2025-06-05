package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arunjeyaprasad/golive/config"
	"github.com/arunjeyaprasad/golive/internal/api/handlers"
	"github.com/arunjeyaprasad/golive/internal/api/middleware"

	"github.com/gorilla/mux"
)

func StartServer() error {
	r := mux.NewRouter()

	// Add middleware
	r.Use(middleware.MuxVars)
	r.Use(middleware.CORS)
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)

	// API routes
	api := r.PathPrefix("/").Subrouter()
	api.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	handlers.RegisterRoutes(api)

	// Create server with timeouts
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", "0.0.0.0", config.DEFAULT_SERVER_PORT),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		slog.Info("Starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	slog.Info("Received shutdown signal", "signal", sig)

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
		return err
	}

	slog.Info("Server shutdown completed")
	return nil
}
