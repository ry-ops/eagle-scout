# Security Vulnerability Fixes

## Scan Date: 2026-02-06

### Vulnerabilities Addressed

**Total vulnerabilities found:** 20 (1 fixable automatically)

### Changes Made

#### 1. Go stdlib Updates (HIGH PRIORITY)
- **CVE-2025-61726** (7.5 H) - Updated Go builder from `golang:1.25-alpine` to `golang:1.25.7-alpine`
- **CVE-2025-61728** (6.5 M) - Fixed by Go 1.25.7 update
- **CVE-2025-61730** (5.3 M) - Fixed by Go 1.25.7 update

#### 2. Base Image Updates
- **CVE-2026-25210** (6.9 M) - Updated from `docker:cli` to `docker:29-cli-alpine3.23` for latest Alpine expat package

#### 3. Remaining Vulnerabilities (Base Image Dependencies)
The following vulnerabilities are in the `docker:29-cli` base image, NOT in eagle-scout's code:
- `github.com/theupdateframework/go-tuf/v2@v2.3.0` (CVE-2026-23992, CVE-2026-23991, CVE-2026-24686)
- `github.com/sigstore/sigstore@v1.10.0` (CVE-2026-24137)
- `github.com/sigstore/rekor@v1.4.3` (CVE-2026-24117, CVE-2026-23831)

**Note:** eagle-scout has NO external Go dependencies - it only uses Go stdlib and wraps the Docker CLI.
These vulnerabilities are bundled in the official Docker CLI image and will be fixed when Docker updates their base image.

### Build Instructions

```bash
# Rebuild the image
docker build -t ryops/eagle-scout:latest .

# Scan the new image
docker scout cves ryops/eagle-scout:latest

# Compare with old image
docker scout compare ryops/eagle-scout:latest --to ryops/eagle-scout:previous
```

### Verification

After rebuild, expected improvements:
- **High severity:** 1 → 0 (100% reduction)
- **Medium severity:** 9 → ~3-4 (pending dependency updates)
- **Total vulnerabilities:** 20 → ~10-12

### Next Steps

1. Update Go dependencies in go.mod
2. Run `go mod tidy`
3. Rebuild Docker image
4. Re-scan with Docker Scout
5. Push updated image to Docker Hub

---

## Scan Date: 2026-02-21

### Upstream-Blocked Vulnerabilities (no action possible)

Scanned `ryops/eagle-scout:1.2.8`. All remaining CVEs are in upstream dependencies we do not own.

#### Alpine 3.23 packages — waiting on Alpine to ship patches

| CVE | Severity | Package | Status |
|-----|----------|---------|--------|
| CVE-2025-15079 | 5.3 M | alpine/curl 8.17.0-r1 | No patched version in Alpine 3.23 repos as of 2026-02-21 |
| CVE-2025-14819 | 5.3 M | alpine/curl 8.17.0-r1 | No patched version in Alpine 3.23 repos as of 2026-02-21 |
| CVE-2025-14524 | 5.3 M | alpine/curl 8.17.0-r1 | No patched version in Alpine 3.23 repos as of 2026-02-21 |
| CVE-2025-15224 | 3.1 L | alpine/curl 8.17.0-r1 | No patched version in Alpine 3.23 repos as of 2026-02-21 |
| CVE-2026-27171 | 2.9 L | alpine/zlib 1.3.1-r2  | No patched version in Alpine 3.23 repos as of 2026-02-21 |

`apk upgrade` already runs in the Dockerfile — it will pick these up automatically once Alpine ships the patches.

#### docker/scout-cli:1.19.2 binary — waiting on Docker to ship a new scout-cli

| CVE | Severity | Package |
|-----|----------|---------|
| CVE-2026-23831 | 5.3 M | github.com/sigstore/rekor 1.4.3 (fix: 1.5.0) |
| CVE-2026-24117 | 5.3 M | github.com/sigstore/rekor 1.4.3 (fix: 1.5.0) |
| CVE-2026-24686 | 4.7 M | github.com/theupdateframework/go-tuf/v2 2.3.0 (fix: 2.4.1) |
| CVE-2026-25934 | 4.3 M | github.com/go-git/go-git/v5 5.16.3 (fix: 5.16.5) |
| CVE-2026-24137 | 5.8 M | github.com/sigstore/sigstore 1.10.3 (fix: 1.10.4) |
| GHSA-mqqf-5wvp-8fh8 | N/A U | github.com/go-chi/chi/v5 5.2.3 (fix: 5.2.4) |
| CVE-2025-68121 | C | stdlib 1.25.6 (fix: 1.25.7) |

#### docker:29.2.1-cli base binary — waiting on Docker CLI update

| CVE | Severity | Package |
|-----|----------|---------|
| CVE-2025-68121 | C | stdlib 1.24.11 (fix: 1.24.13) |
| CVE-2025-61726 | H | stdlib 1.24.11 (fix: 1.24.12) |
| CVE-2025-61728 | M | stdlib 1.24.11 (fix: 1.24.12) |
| CVE-2025-61730 | M | stdlib 1.24.11 (fix: 1.24.12) |

### Action Required

When Alpine ships patched `curl` or `zlib` packages, rebuild and push — `apk upgrade` will pick them up automatically, no Dockerfile changes needed.

When Docker releases a new `scout-cli` version with updated deps, bump the `COPY --from=docker/scout-cli:X.Y.Z` line in the Dockerfile.
