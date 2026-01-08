# Pre-Commit Checklist ✅

Run these commands **before every commit** to ensure CI passes on GitHub.

## Quick Check (30 seconds)

```bash
make fmt && go vet ./... && make build && echo "✓ Ready to push!"
```

If all pass, you're good to commit! If any fail, see detailed checks below.

---

## Detailed Checks

### 1. Format Code ✅

**Command**:
```bash
make fmt
```

**Or manually**:
```bash
gofmt -w .
```

**What it does**: Formats all Go files to standard style

**Why**: GitHub CI will fail if code isn't formatted

---

### 2. Check for Code Issues ✅

**Command**:
```bash
go vet ./...
```

**What it checks**:
- Unused imports
- Unused variables
- Suspicious code patterns
- Common mistakes

**Common fixes**:
```go
// If you have an unused variable that will be used later:
_ = variableName  // Suppress warning

// If you have an unused import:
// Remove it or use it
```

---

### 3. Build All Binaries ✅

**Command**:
```bash
make build
```

**What it does**: Compiles all 5 binaries
- gateway
- datanode
- repair
- objctl
- objbench

**Why**: Ensures code compiles on all platforms

---

## Full Pre-Commit Workflow

```bash
# 1. Make your changes
vim cmd/gateway/main.go

# 2. Format
make fmt

# 3. Check for issues
go vet ./...

# 4. Build
make build

# 5. If you have tests (later)
make test

# 6. Commit
git add .
git commit -m "feat: your change description"

# 7. Push
git push
```

---

## Common Issues & Fixes

### Issue: "code is not formatted"

**Error**:
```
Go code is not formatted:
  cmd/gateway/main.go
```

**Fix**:
```bash
make fmt
# Or
gofmt -w cmd/gateway/main.go
```

---

### Issue: "declared and not used"

**Error**:
```
cmd/repair/main.go:18:2: declared and not used: dbPassword
```

**Fix**:
```go
// Add this line after the variable declaration:
_ = dbPassword  // Will be used when metadata service is implemented
```

---

### Issue: "imported and not used"

**Error**:
```
cmd/datanode/main.go:4:2: "context" imported and not used
```

**Fix**:
```go
// Option 1: Remove the import
// Option 2: Use it in your code
```

---

### Issue: "build failed"

**Error**:
```
undefined: SomeFunction
```

**Fix**:
- Check for typos
- Ensure all imports are correct
- Make sure go.mod is up to date: `go mod tidy`

---

## IDE Integration

### VS Code

Add to `.vscode/settings.json`:
```json
{
  "go.formatTool": "gofmt",
  "editor.formatOnSave": true,
  "[go]": {
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  },
  "go.vetOnSave": "package"
}
```

### GoLand/IntelliJ

1. Preferences → Go → Formatting
2. Enable "Run gofmt on save"
3. Enable "Go vet" inspection

---

## GitHub Actions Status

After pushing, check: https://github.com/mrmushfiq/plinth/actions

You should see:
- ✅ Download dependencies
- ✅ Build all binaries
- ✅ Run go vet
- ✅ Check formatting

If any fail, the commit will show a ❌ instead of ✅

---

## Quick Reference Card

| Command | What | When |
|---------|------|------|
| `make fmt` | Format code | Before every commit |
| `go vet ./...` | Check for issues | Before every commit |
| `make build` | Build binaries | Before every commit |
| `make test` | Run tests | When tests exist |
| `make lint` | Deep linting | Before PR |

---

## Automation Script

Save this as `pre-commit.sh`:

```bash
#!/bin/bash
set -e

echo "🔍 Checking code quality..."

echo "  ✓ Formatting..."
make fmt

echo "  ✓ Running go vet..."
go vet ./...

echo "  ✓ Building..."
make build

echo ""
echo "✅ All checks passed! Safe to commit."
```

Make it executable:
```bash
chmod +x pre-commit.sh
```

Use it:
```bash
./pre-commit.sh && git commit -m "your message"
```

---

## Summary

**Always run before committing**:
```bash
make fmt && go vet ./... && make build
```

**That's it!** These three commands will catch 99% of CI failures before you push.

---

**Pro tip**: Set up a Git pre-commit hook to run these automatically!

See: https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks
