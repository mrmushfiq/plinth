# Contributing to Plinth

Thank you for your interest in contributing to Plinth! This document provides guidelines and information for contributors.

## Code of Conduct

Be respectful, inclusive, and constructive. We're building something useful together.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL 14+
- Basic understanding of distributed systems

### Development Setup

1. **Fork and clone the repository**

```bash
git clone https://github.com/mrmushfiq/plinth.git
cd plinth
```

2. **Install dependencies**

```bash
go mod download
```

3. **Start the development environment**

```bash
make dev-up
```

4. **Run tests**

```bash
make test
```

## Development Workflow

### Branch Naming

- Feature: `feature/your-feature-name`
- Bug fix: `fix/bug-description`
- Documentation: `docs/what-you-changed`
- Refactor: `refactor/what-you-refactored`

### Making Changes

1. Create a branch from `main`
2. Make your changes
3. Write/update tests
4. Ensure all tests pass: `make test`
5. Run linters: `make lint`
6. Format code: `make fmt`
7. Commit with clear messages

### Commit Messages

Follow conventional commits format:

```
type(scope): brief description

Longer explanation if needed.

Fixes #123
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding/updating tests
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Build/tooling changes

Examples:
```
feat(datanode): add xxHash checksum verification

fix(gateway): handle nil pointer in object lookup

docs(api): document multipart upload endpoints
```

### Pull Request Process

1. **Before submitting:**
   - Ensure tests pass locally
   - Update documentation if needed
   - Add tests for new features
   - Run `make lint` and fix any issues

2. **PR Description should include:**
   - What problem does this solve?
   - How does it solve it?
   - Any breaking changes?
   - Screenshots/examples if applicable

3. **PR Review:**
   - Address reviewer feedback
   - Keep PRs focused (one feature/fix per PR)
   - Be patient and respectful

## Testing

### Unit Tests

```bash
# Run all unit tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./internal/placement/...
```

### Integration Tests

```bash
make integration-test
```

### Writing Tests

- Test files should be named `*_test.go`
- Use table-driven tests for multiple cases
- Mock external dependencies
- Aim for >80% coverage on new code

Example:
```go
func TestConsistentHash(t *testing.T) {
    tests := []struct {
        name     string
        key      string
        nodes    []string
        expected string
    }{
        {"basic", "key1", []string{"node1", "node2"}, "node1"},
        // more cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Code Style

### Go Guidelines

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting (automatically via `make fmt`)
- Keep functions small and focused
- Use meaningful variable names
- Add comments for exported functions

### Project-Specific Conventions

1. **Error Handling:**
```go
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

2. **Logging:**
```go
log.Info().
    Str("object_key", key).
    Int("size", size).
    Msg("object stored successfully")
```

3. **Context:**
Always pass `context.Context` as the first parameter for functions that may block or timeout.

## Areas to Contribute

### Good First Issues

Look for issues labeled `good first issue` - these are beginner-friendly tasks.

### Priority Areas

1. **Core Functionality** (Phase 1-2)
   - Data node implementation
   - Replication logic
   - Repair worker

2. **S3 Compatibility** (Phase 3)
   - Additional S3 API endpoints
   - AWS CLI compatibility testing

3. **Unique Features** (Phase 4)
   - Cost tracking enhancements
   - ML-specific optimizations
   - Intelligent tiering

4. **Testing**
   - Integration tests
   - Chaos tests
   - Performance benchmarks

5. **Documentation**
   - Tutorials and guides
   - Architecture documentation
   - API reference

## Documentation

- Update relevant docs in `/docs` directory
- Keep README.md up to date
- Add inline code comments for complex logic
- Include examples in documentation

## Performance Considerations

- Profile before optimizing
- Use benchmarks: `go test -bench=.`
- Consider memory allocation in hot paths
- Document performance characteristics

## Security

### Reporting Vulnerabilities

**Do NOT open public issues for security vulnerabilities.**

Email: security@plinth.dev (or create a GitHub Security Advisory)

Include:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### Security Best Practices

- Never commit secrets/credentials
- Validate all user inputs
- Use parameterized queries for SQL
- Follow secure coding guidelines

## Release Process

(For maintainers)

1. Update version in code
2. Update CHANGELOG.md
3. Tag release: `git tag -a v0.1.0 -m "Release v0.1.0"`
4. Push tag: `git push origin v0.1.0`
5. Create GitHub release with notes

## Communication

- **GitHub Issues**: Bug reports, feature requests
- **GitHub Discussions**: Questions, ideas
- **Pull Requests**: Code contributions

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project README

## Questions?

Don't hesitate to ask! Open a GitHub Discussion or comment on an issue.

---

Thank you for contributing to Plinth! 🚀

