# GitHub Actions Workflows

This directory contains CI/CD workflows for Plinth.

## Active Workflows

### `ci.yml` - Minimal CI (Current)

**Purpose**: Lightweight build verification for early development stage

**Runs on**:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop`

**What it does**:
- ✅ Checks out code
- ✅ Sets up Go 1.21
- ✅ Downloads dependencies
- ✅ Builds all 5 binaries (gateway, datanode, repair, objctl, objbench)
- ✅ Runs `go vet` for basic code quality
- ✅ Checks code formatting with `gofmt`

**Why minimal?**
- Project is in bootstrap/early development phase
- No tests written yet
- Keeps CI green and welcoming for contributors
- Fast feedback (~2 minutes)
- Shows project is actively maintained

## Template Workflows

### `ci-full.yml.template` - Comprehensive CI

**When to activate**: Once you have tests and are ready for production

**To activate**:
```bash
mv .github/workflows/ci-full.yml.template .github/workflows/ci-full.yml
# And optionally remove or rename ci.yml
```

**What it includes**:
- 🧪 **Test Job**: Full test suite with PostgreSQL service
- 🔍 **Lint Job**: golangci-lint for code quality
- 🏗️ **Build Job**: All binaries + artifact upload
- 🐳 **Docker Job**: Build all Docker images

**Additional features**:
- PostgreSQL service for integration tests
- Code coverage reporting to Codecov
- Comprehensive linting
- Docker image building
- Artifact uploads

## Workflow Evolution

### Phase 1: Bootstrap (Current) ✅
**File**: `ci.yml`
- Build verification only
- Fast and simple
- Green builds encourage contributions

### Phase 2: Development (Weeks 2-4)
**Upgrade to**: Add test job
```yaml
# Add when you have tests:
test:
  runs-on: ubuntu-latest
  steps:
    - run: go test ./...
```

### Phase 3: Pre-Production (Weeks 5-8)
**Upgrade to**: `ci-full.yml.template`
- Full test suite
- Linting enforced
- Coverage tracking

### Phase 4: Production (Week 12+)
**Add**: Deployment workflows
- Release automation
- Docker Hub publishing
- Container scanning
- Performance benchmarks

## Adding Tests

When you're ready to add tests:

1. **Create test files**: `*_test.go`
2. **Run locally**: `make test`
3. **Update CI**: Add test job to workflow

Example test job:
```yaml
test:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    - run: go test -v ./...
```

## Badge Status

Add to README.md:
```markdown
[![CI](https://github.com/mrmushfiq/plinth/workflows/CI/badge.svg)](https://github.com/mrmushfiq/plinth/actions)
```

## Troubleshooting

### Build Failing?

1. **Check locally first**:
   ```bash
   go build ./cmd/gateway
   go build ./cmd/datanode
   go vet ./...
   gofmt -l .
   ```

2. **View workflow logs**: 
   - Go to: https://github.com/mrmushfiq/plinth/actions
   - Click on failed run
   - Expand failing step

3. **Common issues**:
   - Missing dependencies: Run `go mod tidy`
   - Formatting: Run `make fmt`
   - Vet errors: Fix issues shown by `go vet`

## Best Practices

1. **Keep CI fast** - Under 5 minutes if possible
2. **Fail fast** - Put quick checks first
3. **Cache dependencies** - Use actions/cache@v3
4. **Matrix testing** - Test multiple Go versions when stable
5. **Protected branches** - Require CI to pass before merge

## Future Enhancements

When ready, consider adding:

- [ ] Matrix builds (Go 1.21, 1.22, 1.23)
- [ ] Integration tests with full cluster
- [ ] Chaos tests in CI
- [ ] Performance regression tests
- [ ] Security scanning (Snyk, Trivy)
- [ ] Dependency updates (Dependabot)
- [ ] Automated releases (GoReleaser)
- [ ] Container publishing (GHCR, Docker Hub)

## Questions?

See: [GitHub Actions Documentation](https://docs.github.com/en/actions)
