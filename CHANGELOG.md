# Changelog

All notable changes to eagle-scout will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- CI/CD pipeline with security gates
- Pull request template
- CHANGELOG.md for tracking changes
- CONTRIBUTING.md guidelines
- SECURITY.md policy
- Branch protection documentation
- SonarCloud integration for code quality analysis
- GitHub Releases with multi-platform binaries (linux/darwin/windows, amd64/arm64)
- GitHub Container Registry (GHCR) publishing alongside Docker Hub
- Multi-architecture Docker images (linux/amd64, linux/arm64)

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

[Unreleased]: https://github.com/ry-ops/eagle-scout/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/ry-ops/eagle-scout/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/ry-ops/eagle-scout/releases/tag/v1.0.0
