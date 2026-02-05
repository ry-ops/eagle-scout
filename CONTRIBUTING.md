# Contributing to eagle-scout

Thank you for your interest in contributing to eagle-scout!

## Development Workflow

### 1. Branch Strategy

```
main (protected)
  └── feat/your-feature    (feature branches)
  └── fix/bug-description  (bug fixes)
  └── chore/updates        (maintenance)
  └── security/cve-fix     (security fixes)
```

- **Never push directly to `main`** - All changes go through pull requests
- Create feature branches from `main`
- Keep branches focused and short-lived

### 2. Before You Start

1. **Check existing issues** - Avoid duplicate work
2. **Research dependencies** - Use latest stable versions
3. **Understand security requirements** - Zero critical/high CVEs policy

### 3. Making Changes

```bash
# Create a feature branch
git checkout main
git pull origin main
git checkout -b feat/your-feature

# Make your changes
# ...

# Run tests locally
go test ./...

# Build and verify
go build -o eagle-scout ./cmd/eagle-scout
./eagle-scout version

# Build Docker image and scan
docker build -t eagle-scout:test .
docker scout cves eagle-scout:test --only-severity critical,high
```

### 4. Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Types:
- `feat` - New feature
- `fix` - Bug fix
- `security` - Security fix
- `chore` - Maintenance
- `docs` - Documentation
- `refactor` - Code refactoring
- `test` - Adding tests

Examples:
```
feat(tools): Add scout_environment tool for environment management

security(deps): Upgrade to Go 1.25 to fix CVE-2025-22871

fix(server): Handle empty image parameter gracefully
```

### 5. Pull Request Process

1. **Update documentation**
   - Add entry to CHANGELOG.md under `[Unreleased]`
   - Update README.md if adding features

2. **Fill out PR template completely**

3. **Ensure CI passes**
   - Build succeeds
   - Tests pass
   - Security scan passes (no new critical/high CVEs)
   - Policy checks pass

4. **Request review**

5. **Address feedback**

### 6. Security Requirements

All contributions must:

- [ ] Introduce no critical or high severity CVEs
- [ ] Maintain non-root container user
- [ ] Not commit secrets or credentials
- [ ] Use dependencies from trusted sources only
- [ ] Pass automated security scanning

### 7. Dependency Updates

When updating dependencies:

1. **Research the update** - Check release notes, breaking changes
2. **Verify security** - Scan for vulnerabilities before and after
3. **Test thoroughly** - Ensure functionality unchanged
4. **Document in CHANGELOG** - Note the update and reason

```bash
# Check for outdated dependencies
go list -u -m all

# Update a specific dependency
go get -u <package>@<version>

# Tidy up
go mod tidy

# Verify
go build ./...
go test ./...
```

## Code Style

- Follow standard Go conventions (`gofmt`)
- Add comments for exported functions
- Keep functions focused and testable
- Handle errors explicitly

## Questions?

Open an issue with the `question` label.
