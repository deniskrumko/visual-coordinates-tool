package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/deniskrumko/visual-coordinates-tool/pkg/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RunServer(ctx context.Context, configPath string) error {
	var config *env.Config

	// Apply config if it's parsed
	if c, err := env.ParseConfig(configPath); err == nil {
		config = c
	}

	indexRouter, err := getIndexRouter(config)
	if err != nil {
		return fmt.Errorf("can't get index router: %w", err)
	}

	recognizeRouter, err := getRecognizeRouter()
	if err != nil {
		return fmt.Errorf("can't get recognize router: %w", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/health"))
	r.Mount("/", indexRouter)
	r.Mount("/recognize", recognizeRouter)

	// Serve static files
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	fileServer(r, "/static", filesDir)

	server := &http.Server{Addr: ":8080", Handler: r}

	// Run server in separate goroutine
	go func() {
		fmt.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Wait for shutdown
	<-ctx.Done()

	// Set timeout to finish active requests
	fmt.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shutdown error: %v", err)
	}

	fmt.Println("Server gracefully stopped.")
	return nil
}
