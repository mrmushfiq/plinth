package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Configuration from environment
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "plinth")
	dbUser := getEnv("DB_USER", "plinth")
	dbPassword := getEnv("DB_PASSWORD", "plinth_dev_password")
	_ = dbPassword // Will be used when metadata service is implemented
	repairInterval := getEnvDuration("REPAIR_INTERVAL", 60*time.Second)
	scrubInterval := getEnvDuration("SCRUB_INTERVAL", 300*time.Second)

	log.Printf("Starting Plinth Repair Worker")
	log.Printf("Database: %s@%s:%s/%s", dbUser, dbHost, dbPort, dbName)
	log.Printf("Repair interval: %s", repairInterval)
	log.Printf("Scrub interval: %s", scrubInterval)

	// TODO: Initialize metadata service
	// TODO: Initialize data node clients

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down repair worker...")
		cancel()
	}()

	// Repair loop
	repairTicker := time.NewTicker(repairInterval)
	defer repairTicker.Stop()

	// Scrub loop
	scrubTicker := time.NewTicker(scrubInterval)
	defer scrubTicker.Stop()

	log.Println("Repair worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Repair worker stopped")
			return
		case <-repairTicker.C:
			log.Println("Running repair cycle...")
			// TODO: Implement repair logic
			// - Find under-replicated objects
			// - Rebuild from healthy replicas
		case <-scrubTicker.C:
			log.Println("Running scrub cycle...")
			// TODO: Implement scrub logic
			// - Verify checksums
			// - Detect corruption
		}
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}
