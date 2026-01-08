package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	// Configuration from environment
	nodeID := getEnv("NODE_ID", "node1")
	dataDir := getEnv("DATA_DIR", "./data")
	grpcPort := getEnv("GRPC_PORT", "50051")

	log.Printf("Starting Plinth Data Node: %s", nodeID)
	log.Printf("Data directory: %s", dataDir)
	log.Printf("gRPC port: %s", grpcPort)

	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// TODO: Initialize storage service
	// TODO: Register gRPC services

	// Setup gRPC server
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// TODO: Register storage service
	// pb.RegisterStorageServiceServer(grpcServer, storageService)

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down data node...")
		grpcServer.GracefulStop()
	}()

	log.Printf("Data node %s listening on :%s", nodeID, grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
