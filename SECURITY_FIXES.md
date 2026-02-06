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
