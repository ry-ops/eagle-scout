# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.1.x   | :white_check_mark: |
| 1.0.x   | :x:                |

## Security Standards

eagle-scout maintains the following security standards:

### Container Security
- **Non-root execution** - Container runs as user `scout` (UID 1001)
- **Minimal base image** - Uses `docker:cli` with only required tools
- **No secrets in image** - All credentials passed at runtime
- **Supply chain attestations** - SBOM and provenance included

### Vulnerability Policy
- **Critical CVEs** - Must be fixed before merge
- **High CVEs** - Must be fixed before merge (if fix available)
- **Medium/Low CVEs** - Tracked and fixed in regular updates

### CI/CD Security Gates
All pull requests must pass:
1. Docker Scout CVE scan (critical/high with fixes = fail)
2. Non-root user check
3. SBOM generation
4. Provenance attestation

## Reporting a Vulnerability

**DO NOT** open a public issue for security vulnerabilities.

### Private Disclosure

1. Email: security@ry-ops.dev (or open a private security advisory on GitHub)
2. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

### Response Timeline

- **Acknowledgment** - Within 48 hours
- **Initial Assessment** - Within 7 days
- **Fix Timeline** - Based on severity:
  - Critical: 24-48 hours
  - High: 7 days
  - Medium: 30 days
  - Low: Next release

## Security Scanning

We use Docker Scout for continuous security monitoring:

```bash
# Scan for vulnerabilities
docker scout cves ryops/eagle-scout:latest

# Quick overview
docker scout quickview ryops/eagle-scout:latest

# Check policy compliance
docker scout policy ryops/eagle-scout:latest
```

## Known Issues

### Upstream Dependencies

Some vulnerabilities may exist in upstream dependencies that we cannot directly fix:

| Component | Issue | Status | Tracking |
|-----------|-------|--------|----------|
| docker:cli | Go 1.24.11 stdlib CVE | Waiting on Docker | [docker/docker#xxx] |

We monitor upstream fixes and update promptly when available.

## Security Changelog

### 2026-02-05 (v1.1.0)
- Fixed 7 Go stdlib CVEs by upgrading to Go 1.25
- Added SBOM attestation
- Added max-mode provenance attestation

### 2026-02-05 (v1.0.0)
- Initial release with non-root user
- Multi-stage build for minimal image
