# Changelog

All notable changes to eagle-scout will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.2.4] - 2026-02-07

### Security
- Pinned base image to `docker:29.2.1-cli` (Docker CLI built with Go 1.25.6+)
  - CVE-2025-61726 (high 7.5) - FIXED
  - CVE-2025-61728 (medium 6.5) - FIXED
  - CVE-2025-61730 (medium 5.3) - FIXED
- Upgraded Alpine expat to 2.7.4-r0
  - CVE-2026-25210 (medium 6.9) - FIXED
- Pinned Docker Scout CLI to v1.19.0 for reproducibility
- Remaining 5 medium CVEs in Scout CLI dependencies (go-tuf, sigstore, rekor) require upstream update

## [1.2.3] - 2026-02-07

### Changed
- Version bump to consolidate v1.2.1 and v1.2.2 fixes into a clean release

## [1.2.2] - 2026-02-07

### Fixed
- Added Docker Scout CLI plugin to container image — scout commands now work inside the container
- Previously, `docker scout` was not available inside the image, causing eagle-scout to fail on startup

### Added
- Documented native multi-arch support (linux/amd64, linux/arm64) in README
- Docker images now explicitly noted as working natively on Intel and Apple Silicon

## [1.2.1] - 2026-02-07

### Added
- Documented native multi-arch support (linux/amd64, linux/arm64) in README

### Fixed
- Clarified multi-arch image availability for ARM-based systems (e.g., Apple M1/M2/M3)

## [1.2.0] - 2026-02-07

### Fixed
- Policy checks (secrets detection, attack surface) now properly fail the build instead of silently passing
- CVE gate now blocks on ALL critical/high CVEs, not just those with available fixes
- Multi-arch builds (linux/amd64, linux/arm64) now verified in CI on every push/PR
- Publish step now builds multi-arch images instead of single-arch
- Release workflow artifact download no longer picks up Docker buildx cache artifacts
- Security scan builds without Docker layer cache to prevent stale CVE false passes

### Changed
- Branch protection now enforces admin compliance (no bypass)
- CI/CD pipeline with security gates
- SonarCloud integration for code quality analysis
- GitHub Releases with multi-platform binaries (linux/darwin/windows, amd64/arm64)
- GitHub Container Registry (GHCR) publishing alongside Docker Hub

## [1.1.0] - 2026-02-05

### Added
- `scout_environment` tool - Manage environments (list/set) for policy evaluation
- `scout_cache` tool - Manage local cache (df/prune)
- `scout_enroll` tool - Enroll organization with Docker Scout
- `scout_watch` tool - Enable/disable continuous monitoring
- `scout_vex` now supports `list` action

### Security
- Upgraded from Go 1.22 to Go 1.25 to fix critical stdlib CVEs:
  - CVE-2025-22871 (critical) - FIXED
  - CVE-2025-61729 (high) - FIXED
  - CVE-2025-61726 (high) - FIXED
  - CVE-2025-61725 (high) - FIXED
  - CVE-2025-61723 (high) - FIXED
  - CVE-2025-58188 (high) - FIXED
  - CVE-2025-58187 (high) - FIXED
- Added SBOM attestation to published images
- Added max-mode provenance attestation

### Changed
- Build image updated to `golang:1.25-alpine`
- Total tools increased from 10 to 15

## [1.0.0] - 2026-02-05

### Added
- Initial release with 10 Docker Scout tools:
  - `scout_cves` - Scan images for CVEs
  - `scout_quickview` - Quick security overview
  - `scout_compare` - Compare two images
  - `scout_sbom` - Generate SBOM
  - `scout_recommendations` - Base image recommendations
  - `scout_policy` - Policy evaluation
  - `scout_attestation` - Manage attestations
  - `scout_repo` - Repository management
  - `scout_vex` - VEX statement management
  - `scout_version` - Version info
- Animated compass logo
- Docker image with non-root user
- MCP server over stdio transport

### Security
- Runs as non-root user (scout, UID 1001)
- Multi-stage build to minimize image size

---

[Unreleased]: https://github.com/ry-ops/eagle-scout/compare/v1.2.4...HEAD
[1.2.4]: https://github.com/ry-ops/eagle-scout/compare/v1.2.3...v1.2.4
[1.2.3]: https://github.com/ry-ops/eagle-scout/compare/v1.2.2...v1.2.3
[1.2.2]: https://github.com/ry-ops/eagle-scout/compare/v1.2.1...v1.2.2
[1.2.1]: https://github.com/ry-ops/eagle-scout/compare/v1.2.0...v1.2.1
[1.2.0]: https://github.com/ry-ops/eagle-scout/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/ry-ops/eagle-scout/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/ry-ops/eagle-scout/releases/tag/v1.0.0
