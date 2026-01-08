# API Package

This package contains the HTTP API implementation for Plinth's S3-compatible gateway.

## Structure

```
api/
├── router.go      # Main router setup and route definitions
├── handlers.go    # S3 API handler implementations
├── middleware.go  # Custom middleware (auth, metrics, etc.)
└── README.md      # This file
```

## Components

### Router (`router.go`)

Sets up the Gin router with all routes and middleware:
- Health checks and metrics endpoints
- Admin API (cluster status, costs, repair)
- S3 API routes (bucket and object operations)

### Handlers (`handlers.go`)

Implements S3-compatible API handlers:
- **Bucket operations**: Create, Delete, Head, List
- **Object operations**: Put, Get, Delete, Head
- **Multipart uploads**: Initiate, UploadPart, Complete, Abort, List
- **Error responses**: S3-compatible error codes

### Middleware (`middleware.go`)

Custom middleware functions:
- `RequestIDMiddleware`: Adds unique request ID to each request
- `CORSMiddleware`: Handles CORS for S3 compatibility
- `MetricsMiddleware`: Collects Prometheus metrics
- `AuthMiddleware`: AWS SigV4 authentication (to be implemented)
- `LoggingMiddleware`: Structured request logging

## Usage

```go
import "github.com/mushfiq/plinth/internal/api"

// Create gateway with dependencies
gateway := api.NewGateway()

// Setup router
router := api.SetupRouter(gateway, "development")

// Start server
router.Run(":9000")
```

## S3 API Coverage

### Implemented (Stubs)
- ✅ Bucket operations (Create, Delete, Head, List)
- ✅ Object operations (Put, Get, Delete, Head)
- ✅ Multipart uploads (all operations)
- ✅ List objects

### To Be Implemented
- [ ] AWS SigV4 authentication
- [ ] Pre-signed URLs
- [ ] Object versioning
- [ ] Object tagging
- [ ] Bucket policies
- [ ] CORS configuration
- [ ] Lifecycle policies

## Adding New Handlers

1. Add handler method to `Gateway` struct in `handlers.go`:
```go
func (g *Gateway) MyNewHandler(c *gin.Context) {
    // Implementation
}
```

2. Register route in `router.go`:
```go
router.GET("/my-endpoint", gateway.MyNewHandler)
```

3. Add tests in `handlers_test.go`

## Testing

```bash
# Run API tests
go test ./internal/api/...

# With coverage
go test -cover ./internal/api/...
```

## S3 Error Codes

Common error codes implemented:
- `NoSuchBucket`: Bucket doesn't exist
- `NoSuchKey`: Object doesn't exist
- `BucketAlreadyExists`: Bucket name conflict
- `InvalidBucketName`: Invalid bucket name
- `InvalidArgument`: Invalid parameter
- `InternalError`: Server error
- `AccessDenied`: Authentication/authorization failure

See `handlers.go` for full list.

## Dependencies

- `github.com/gin-gonic/gin` - Web framework
- `internal/metadata` - Metadata service interface
- `internal/placement` - Placement controller interface
- `internal/quorum` - Quorum write/read logic

## Future Improvements

1. **Performance**
   - Connection pooling for data nodes
   - Response streaming for large objects
   - Caching for metadata lookups

2. **Features**
   - Complete AWS SigV4 implementation
   - Pre-signed URL generation and validation
   - Object versioning support
   - Server-side encryption

3. **Observability**
   - Distributed tracing (OpenTelemetry)
   - Detailed request metrics
   - Audit logging

4. **Security**
   - Rate limiting per client
   - Bucket policies
   - IAM-like access control

