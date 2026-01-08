package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrmushfiq/plinth/internal/api"
)

func main() {
	// Configuration from environment
	port := getEnv("HTTP_PORT", "9000")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "plinth")
	dbUser := getEnv("DB_USER", "plinth")
	dbPassword := getEnv("DB_PASSWORD", "plinth_dev_password")
	environment := getEnv("ENVIRONMENT", "development")

	log.Printf("Starting Plinth Gateway on port %s", port)
	log.Printf("Database: %s@%s:%s/%s", dbUser, dbHost, dbPort, dbName)
	log.Printf("Environment: %s", environment)

	// TODO: Initialize metadata service (use dbPassword)
	// TODO: Initialize placement service
	// TODO: Initialize data node clients
	_ = dbPassword // Will be used when metadata service is implemented

	// Create gateway with dependencies
	gateway := api.NewGateway()

	// Setup Gin router
	router := api.SetupRouter(gateway, environment)

	// Create HTTP server
	server := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Graceful shutdown handler
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down gateway...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	log.Printf("Gateway listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
