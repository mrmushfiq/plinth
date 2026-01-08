# Plinth

> Cost-aware, self-healing object storage built for ML workloads and on-prem clusters

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![CI](https://github.com/mrmushfiq/plinth/workflows/CI/badge.svg)](https://github.com/mrmushfiq/plinth/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/mrmushfiq/plinth)](https://goreportcard.com/report/github.com/mrmushfiq/plinth)

**Plinth** is an open-source, S3-compatible distributed object storage system designed for AI/ML workloads and on-premises clusters. Unlike traditional object stores, Plinth provides built-in cost tracking, intelligent tiering, and transparent self-healing capabilities.

## 🎯 Why Plinth?

- **Cost Visibility**: Track storage costs per bucket/dataset/experiment
- **ML-Optimized**: Built for ML access patterns (large sequential reads, dataset versioning)
- **Self-Healing**: Automatic detection and repair of under-replicated objects
- **Simple Operations**: Easier to deploy and debug than MinIO or Ceph
- **S3-Compatible**: Drop-in replacement for basic S3 workflows

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/mrmushfiq/plinth.git
cd plinth

# Start local cluster with docker-compose
docker-compose up -d

# Use with AWS CLI
aws --endpoint-url=http://localhost:9000 s3 mb s3://my-bucket
aws --endpoint-url=http://localhost:9000 s3 cp file.txt s3://my-bucket/
```

## 📐 Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     API Gateway                          │
│  (S3-compatible HTTP, Auth, Request Routing)            │
└────────────┬────────────────────────────────────────────┘
             │
             ├──────────────┬──────────────────────────────┐
             │              │                              │
┌────────────▼────┐  ┌──────▼──────┐         ┌───────────▼────┐
│  Metadata DB    │  │   Placement  │         │  Data Nodes    │
│  (PostgreSQL)   │  │   Controller │         │  (3+ nodes)    │
│                 │  │              │         │                │
│ • Object keys   │  │ • Consistent │         │ • Local disk   │
│ • Versions      │  │   hashing    │         │ • Checksums    │
│ • Placement     │  │ • Rebalance  │         │ • gRPC API     │
│ • Costs         │  │              │         │                │
└─────────────────┘  └──────────────┘         └────────────────┘
                              │
                     ┌────────▼─────────┐
                     │  Repair Worker   │
                     │                  │
                     │ • Detect issues  │
                     │ • Rebuild copies │
                     │ • Scrub data     │
                     └──────────────────┘
```

## ✨ Features

### Current (v0.1.0-alpha)
- [ ] Basic PUT/GET/DELETE operations
- [ ] Multi-node replication (configurable RF)
- [ ] Checksum verification (xxHash)
- [ ] PostgreSQL metadata storage
- [ ] Consistent hashing for placement

### Planned
- [ ] Automatic repair and self-healing
- [ ] Cost tracking per bucket
- [ ] Multipart uploads
- [ ] Intelligent tiering (hot/warm/cold)
- [ ] ML-friendly features (dataset manifests, batch operations)
- [ ] Pre-signed URLs
- [ ] Object versioning

See [plinth_project_plan.md](./plinth_project_plan.md) for the full roadmap.

## 🛠️ Development

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 14+

### Build from Source

```bash
# Build all binaries
make build

# Run tests
make test

# Run integration tests
make integration-test

# Start development cluster
make dev-up
```

### Project Structure

```
plinth/
├── cmd/                   # Binary entry points
│   ├── gateway/           # API gateway
│   ├── datanode/          # Storage node
│   ├── repair/            # Repair worker
│   ├── objctl/            # Admin CLI
│   └── objbench/          # Benchmark tool
├── internal/              # Private packages
│   ├── api/               # HTTP handlers
│   ├── metadata/          # PostgreSQL repo
│   ├── placement/         # Consistent hashing
│   ├── quorum/            # Quorum writes
│   ├── checksum/          # Integrity checks
│   ├── tiering/           # Storage tiering
│   └── cost/              # Cost tracking
├── pkg/                   # Public APIs
├── deploy/                # Deployment configs
├── docs/                  # Documentation
└── chaos/                 # Chaos tests
```

## 📊 Benchmarks

Coming soon! We'll compare against MinIO and Ceph.

## 🧪 Testing

```bash
# Unit tests
make test

# Integration tests
make integration-test

# Chaos tests
make chaos-kill-node
make chaos-network-partition
make chaos-disk-full
```

## 📚 Documentation

- [Quick Start Guide](docs/quickstart.md) (coming soon)
- [Architecture Overview](docs/architecture.md) (coming soon)
- [API Reference](docs/api/s3-compatibility.md) (coming soon)
- [Operations Guide](docs/operations/deployment.md) (coming soon)

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Inspired by:
- [Amazon S3](https://aws.amazon.com/s3/)
- [MinIO](https://min.io/)
- [Ceph](https://ceph.io/)
- The Dynamo paper and consistent hashing principles

## 📬 Contact

- GitHub Issues: [github.com/mrmushfiq/plinth/issues](https://github.com/mrmushfiq/plinth/issues)
- Project Maintainer: [@mrmushfiq](https://github.com/mrmushfiq)

---

**Status**: 🚧 Early Development - Not ready for production use

