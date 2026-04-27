package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/queue/internal/config"
	"github.com/user/queue/internal/database"
	"github.com/user/queue/internal/router"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("application error: %v", err)
	}
}

func run() error {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database pool
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbPool, err := database.NewPool(ctx, cfg)
	if err != nil {
		return fmt.Errorf("init db pool: %w", err)
	}

	defer dbPool.Close()

	r := router.Setup()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Server error channel
	errCh := make(chan error, 1)

	// Start server
	go func() {
		log.Printf("server starting on localhost:%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("server exited gracefully")
	return nil
}
