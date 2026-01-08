.PHONY: help build test integration-test clean dev-up dev-down docker-build lint

# Variables
BINARY_DIR := bin
GATEWAY_BINARY := $(BINARY_DIR)/gateway
DATANODE_BINARY := $(BINARY_DIR)/datanode
REPAIR_BINARY := $(BINARY_DIR)/repair
OBJCTL_BINARY := $(BINARY_DIR)/objctl
OBJBENCH_BINARY := $(BINARY_DIR)/objbench

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod
GOVET := $(GOCMD) vet
GOFMT := gofmt

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all binaries
	@echo "Building binaries..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(GATEWAY_BINARY) ./cmd/gateway
	$(GOBUILD) -o $(DATANODE_BINARY) ./cmd/datanode
	$(GOBUILD) -o $(REPAIR_BINARY) ./cmd/repair
	$(GOBUILD) -o $(OBJCTL_BINARY) ./cmd/objctl
	$(GOBUILD) -o $(OBJBENCH_BINARY) ./cmd/objbench
	@echo "Build complete!"

test: ## Run unit tests
	@echo "Running unit tests..."
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

integration-test: ## Run integration tests
	@echo "Running integration tests..."
	$(GOTEST) -v -tags=integration ./...

lint: ## Run linters
	@echo "Running linters..."
	$(GOVET) ./...
	@echo "Checking formatting..."
	@test -z "$$($(GOFMT) -l .)" || (echo "Code needs formatting. Run 'make fmt'" && exit 1)

fmt: ## Format Go code
	@echo "Formatting code..."
	$(GOFMT) -w .

mod-tidy: ## Tidy Go modules
	$(GOMOD) tidy

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf $(BINARY_DIR)
	rm -f coverage.txt coverage.html
	@echo "Clean complete!"

dev-up: ## Start development environment
	@echo "Starting development cluster..."
	docker-compose up -d
	@echo "Cluster started! Gateway: http://localhost:9000"

dev-down: ## Stop development environment
	@echo "Stopping development cluster..."
	docker-compose down -v

dev-logs: ## Show logs from development cluster
	docker-compose logs -f

docker-build: ## Build Docker images
	@echo "Building Docker images..."
	docker build -t plinth/gateway:dev -f deploy/Dockerfile.gateway .
	docker build -t plinth/datanode:dev -f deploy/Dockerfile.datanode .
	docker build -t plinth/repair:dev -f deploy/Dockerfile.repair .

# Chaos testing targets
chaos-kill-node: ## Kill a random data node during operation
	@echo "Running chaos test: kill node"
	./chaos/kill_node.sh

chaos-network-partition: ## Simulate network partition
	@echo "Running chaos test: network partition"
	./chaos/network_partition.sh

chaos-disk-full: ## Simulate disk full scenario
	@echo "Running chaos test: disk full"
	./chaos/disk_full.sh

chaos-corrupt-replica: ## Corrupt a replica on disk
	@echo "Running chaos test: corrupt replica"
	./chaos/corrupt_replica.sh

# Benchmark targets
bench-put: ## Run PUT benchmark
	$(OBJBENCH_BINARY) put --concurrency=50 --size=1MB --count=10000

bench-get: ## Run GET benchmark
	$(OBJBENCH_BINARY) get --concurrency=100 --size=10MB --count=5000

bench-all: bench-put bench-get ## Run all benchmarks

