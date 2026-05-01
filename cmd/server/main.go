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
	// Load configuration (e.g env)
	cfg := config.LoadConfig()

	// Set up global signal notification for shutdown
	// SIGINT - CTRL+C
	// SIGTERM - termination signal
	// It notifies the context channel when a shutdown signal is received
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop() // stops the signal notification when the function exits

	// Initialize database pool
	dbPool, err := database.NewPool(ctx, cfg)
	if err != nil {
		return fmt.Errorf("error while initializing db pool: %w", err)
	}
	defer dbPool.Close() // stops the database pool when the function exits

	r := router.Setup()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.PORT),
		Handler: r,
	}

	// Server error channel
	errCh := make(chan error, 1)

	// Start server via goroutine because ListenAndServe blocks until the server is closed
	go func() {
		log.Printf("server starting on localhost:%s", srv.Addr)

		// ListenAndServe returns an error if the server is closed unexpectedly
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// send errors to error channel which are not server closed errors
			errCh <- err
		}
	}()

	// select blocks until one of the cases (shutdown signal or server error) is ready to execute
	// if the context is done, the shutdown signal is received
	// if an error is received, the server error is returned
	// if neither is ready, select blocks indefinitely
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
