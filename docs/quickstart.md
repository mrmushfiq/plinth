# Quick Start Guide

Get Plinth up and running in 5 minutes.

## Prerequisites

- Docker and Docker Compose
- AWS CLI (for testing S3 compatibility)

## Option 1: Docker Compose (Recommended)

### 1. Clone the Repository

```bash
git clone https://github.com/mushfiq/plinth.git
cd plinth
```

### 2. Start the Cluster

```bash
docker-compose up -d
```

This will start:
- 3 data nodes (ports 50051-50053)
- 1 API gateway (port 9000)
- 1 PostgreSQL database
- 1 repair worker

### 3. Verify the Cluster

```bash
# Check cluster health
curl http://localhost:9000/health

# Check container status
docker-compose ps
```

### 4. Configure AWS CLI

Add a profile for Plinth:

```bash
aws configure --profile plinth
# AWS Access Key ID: minioadmin
# AWS Secret Access Key: minioadmin
# Default region name: us-east-1
# Default output format: json
```

### 5. Test Basic Operations

```bash
# Create a bucket
aws --profile plinth --endpoint-url=http://localhost:9000 s3 mb s3://test-bucket

# Upload a file
echo "Hello, Plinth!" > test.txt
aws --profile plinth --endpoint-url=http://localhost:9000 s3 cp test.txt s3://test-bucket/

# List objects
aws --profile plinth --endpoint-url=http://localhost:9000 s3 ls s3://test-bucket/

# Download the file
aws --profile plinth --endpoint-url=http://localhost:9000 s3 cp s3://test-bucket/test.txt downloaded.txt

# Verify
cat downloaded.txt
```

## Option 2: Build from Source

### 1. Install Go 1.21+

```bash
# macOS
brew install go

# Linux
# Download from https://golang.org/dl/
```

### 2. Clone and Build

```bash
git clone https://github.com/mushfiq/plinth.git
cd plinth

# Download dependencies
go mod download

# Build all binaries
make build
```

### 3. Start PostgreSQL

```bash
# Using Docker
docker run -d \
  --name plinth-postgres \
  -e POSTGRES_DB=plinth \
  -e POSTGRES_USER=plinth \
  -e POSTGRES_PASSWORD=plinth_dev_password \
  -p 5432:5432 \
  postgres:14-alpine

# Run init script
psql -h localhost -U plinth -d plinth < deploy/sql/init.sql
```

### 4. Start Data Nodes

```bash
# Terminal 1
./bin/datanode --node-id=node1 --port=50051 --data-dir=./data/node1

# Terminal 2
./bin/datanode --node-id=node2 --port=50052 --data-dir=./data/node2

# Terminal 3
./bin/datanode --node-id=node3 --port=50053 --data-dir=./data/node3
```

### 5. Start Gateway

```bash
# Terminal 4
./bin/gateway \
  --port=9000 \
  --db-host=localhost \
  --db-port=5432 \
  --db-name=plinth \
  --data-nodes=localhost:50051,localhost:50052,localhost:50053
```

### 6. Start Repair Worker

```bash
# Terminal 5
./bin/repair \
  --db-host=localhost \
  --db-port=5432 \
  --data-nodes=localhost:50051,localhost:50052,localhost:50053
```

## Using the Admin CLI

### Check Cluster Status

```bash
./bin/objctl cluster status
```

### List Nodes

```bash
./bin/objctl nodes list
```

### Check Object Details

```bash
./bin/objctl object stat test-bucket/test.txt
```

### View Costs

```bash
./bin/objctl costs bucket test-bucket
```

## Running Tests

```bash
# Unit tests
make test

# Integration tests
make integration-test

# With coverage
make test-coverage
```

## Troubleshooting

### Gateway not starting

Check if PostgreSQL is running:
```bash
docker ps | grep postgres
```

Check gateway logs:
```bash
docker-compose logs gateway
```

### Data node connection issues

Verify data nodes are healthy:
```bash
docker-compose ps | grep datanode
```

Check data node logs:
```bash
docker-compose logs datanode1
```

### Port already in use

Stop existing services:
```bash
docker-compose down
```

Or change ports in `docker-compose.yml`.

## Next Steps

- Read the [Architecture Overview](architecture.md)
- Learn about [Consistency Guarantees](consistency.md)
- Deploy to [Kubernetes](operations/deployment.md)
- Enable [Cost Tracking](ml-guide.md#cost-tracking)

## Getting Help

- GitHub Issues: https://github.com/mushfiq/plinth/issues
- Discussions: https://github.com/mushfiq/plinth/discussions
- Documentation: https://github.com/mushfiq/plinth/tree/main/docs

