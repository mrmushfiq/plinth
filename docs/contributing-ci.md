# CI/CD Guidelines for Contributors

This document explains how Continuous Integration works for Plinth and what to expect when submitting Pull Requests.

## Current CI Setup

Plinth uses **GitHub Actions** for continuous integration. We currently run a **minimal, fast CI** focused on build verification.

## What Runs on Your PR

When you submit a pull request, the following checks will run automatically:

### 1. Build Verification ⚡
- Compiles all 5 binaries
- Ensures code compiles on Linux/amd64
- **Runtime**: ~2 minutes

### 2. Code Quality Checks 🔍
- `go vet`: Catches common errors
- `gofmt`: Verifies code formatting

### 3. Dependency Check 📦
- Ensures `go.mod` and `go.sum` are valid
- Downloads and caches dependencies

## Before Submitting a PR

Run these commands locally to ensure CI passes:

```bash
# Build all binaries
make build

# Run go vet
go vet ./...

# Check formatting
gofmt -l .

# If any files are listed, format them:
make fmt

# Verify everything
go build ./cmd/...
```

## CI Status

You'll see these checks on your PR:

- ✅ **Build / Build** - All binaries compiled successfully
- ✅ **Build / go vet** - Code passed static analysis  
- ✅ **Build / formatting** - Code is properly formatted

All must pass before the PR can be merged.

## What If CI Fails?

### Build Failure

**Error**: `build failed: undefined reference`

**Fix**: Missing import or typo
```bash
go build ./cmd/gateway  # Test locally
```

### Vet Failure

**Error**: `go vet: suspicious assignment`

**Fix**: Address the issue reported by `go vet`
```bash
go vet ./...
```

### Formatting Failure

**Error**: `code is not formatted`

**Fix**: Run `gofmt` or use your editor's auto-format
```bash
make fmt
# Or
gofmt -w .
```

## Why Minimal CI?

Plinth is in early development (bootstrap phase). We use minimal CI to:

1. **Encourage contributions** - Fast, simple checks
2. **Stay green** - No failing tests to worry about (yet!)
3. **Focus on building** - Not testing/deployment (yet)
4. **Quick feedback** - Results in ~2 minutes

## What's Coming

As the project matures, we'll add:

### Phase 2: Tests (Weeks 2-4)
```yaml
# Will add:
- run: go test ./...
```

### Phase 3: Comprehensive CI (Weeks 5+)
- Full test suite with PostgreSQL
- Code coverage tracking
- golangci-lint for deeper analysis
- Integration tests

### Phase 4: Advanced CI (Production)
- Docker image building
- Performance benchmarks
- Security scanning
- Automated releases

## IDE Integration

### VS Code

Install the Go extension and it will:
- Auto-format on save (gofmt)
- Show vet errors inline
- Run tests in IDE

**Settings**:
```json
{
  "go.formatTool": "gofmt",
  "editor.formatOnSave": true,
  "[go]": {
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  }
}
```

### GoLand / IntelliJ

- Enable "gofmt" in Preferences → Go → Formatting
- Enable "Run gofmt on save"
- Enable "Go vet" inspection

## Local Development Workflow

Recommended workflow to avoid CI failures:

```bash
# 1. Make your changes
vim cmd/gateway/main.go

# 2. Format code
make fmt

# 3. Build to check compilation
make build

# 4. Run vet
go vet ./...

# 5. When tests exist, run them
make test

# 6. Commit and push
git add .
git commit -m "feat: add new feature"
git push
```

## Questions?

- Check `.github/workflows/README.md` for workflow details
- View workflow runs: https://github.com/mrmushfiq/plinth/actions
- Ask in GitHub Discussions

## Summary

✅ **Do this before submitting PR**:
```bash
make fmt      # Format code
make build    # Verify builds
go vet ./...  # Check for issues
```

❌ **Don't worry about**:
- Tests (none yet!)
- Linting (minimal for now)
- Docker builds (not in CI)
- Code coverage (coming later)

Keep it simple, keep it building, and we'll expand as the project grows! 🚀
