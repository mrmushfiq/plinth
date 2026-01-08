# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project bootstrap with complete structure
- Go module setup with dependencies
- Docker and docker-compose configuration for development
- PostgreSQL schema with complete metadata model
- Gin-based HTTP API gateway with S3-compatible routes
- Comprehensive middleware (CORS, request ID, metrics, logging)
- API handler stubs for all S3 operations:
  - Bucket operations (Create, Delete, Head, List)
  - Object operations (Put, Get, Delete, Head)
  - Multipart uploads (Initiate, UploadPart, Complete, Abort, List)
- Admin API endpoints:
  - `/admin/cluster/status` - Cluster health
  - `/admin/nodes` - Node status
  - `/admin/costs/by-bucket` - Cost tracking
  - `/admin/costs/top-objects` - Top objects by cost
  - `/admin/repair/status` - Repair worker status
- Health check endpoint (`/health`)
- Metrics endpoint (`/metrics`) for Prometheus
- Binary stubs for all services (gateway, datanode, repair, objctl, objbench)
- Core package interfaces:
  - `metadata.Service` - Metadata operations
  - `placement.Controller` - Node placement
  - `quorum.Writer/Reader` - Quorum operations
  - `checksum.Calculator` - Checksum verification
- Documentation:
  - README.md with project overview
  - CONTRIBUTING.md with development guidelines
  - docs/quickstart.md for users
  - docs/architecture.md with system design
  - internal/api/README.md for API package
- CI/CD pipeline with GitHub Actions
- Makefile with development commands

### Changed
- Switched from stdlib `net/http` mux to Gin framework for better performance and features
- Simplified CI workflow to minimal build verification (moved full CI to template for later use)

### Technical Details

**Framework Choice: Gin**
- Better routing for complex S3 API (path params, query params)
- Built-in middleware support
- Superior performance with httprouter
- Excellent XML/JSON handling for S3 responses
- Active community and good documentation

**Architecture:**
- Clean separation: router → handlers → middleware
- Structured error responses with S3 error codes
- Request ID tracking for observability
- CORS support for browser clients
- Graceful shutdown handling

## [0.1.0-alpha] - 2026-01-08

### Added
- Initial release (bootstrap)
- Project structure and scaffolding
- All services compile successfully
- Gateway responds to health checks
- Complete S3 route structure (handlers to be implemented)

---

## Development Status

- ✅ Phase 0: Bootstrap Complete
- 🔨 Phase 1, Week 1: In Progress (Data Node Implementation)
- ⏳ Remaining: 11 weeks according to project plan

