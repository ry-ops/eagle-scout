# eagle-scout

<p align="center">
  <img src="assets/eagle-scout-logo.svg" alt="eagle-scout logo" width="400"/>
</p>

**MCP Server for Docker Scout** - Container security scanning via Model Context Protocol.

[![Docker Scout Grade](https://img.shields.io/badge/Docker_Scout-Grade_A-success)](https://hub.docker.com/r/ryops/eagle-scout)
[![Security](https://img.shields.io/badge/vulnerabilities-14-yellow)](https://github.com/ry-ops/eagle-scout/blob/main/SECURITY_FIXES.md)
[![Version](https://img.shields.io/badge/version-1.1.1-blue)](https://github.com/ry-ops/eagle-scout/releases/tag/v1.1.1)

Part of the [ry-ops](https://github.com/ry-ops) fabric ecosystem.

## Features

- **CVE Scanning** - Scan container images for vulnerabilities
- **Quick Overview** - Get instant security summaries
- **Image Comparison** - Diff two images for security changes
- **SBOM Generation** - Software Bill of Materials in SPDX/CycloneDX
- **Recommendations** - Base image update suggestions
- **Policy Evaluation** - Check images against security policies
- **Attestations** - Manage supply chain attestations
- **VEX Management** - Vulnerability Exploitability eXchange

## Prerequisites

- Docker Desktop 4.17+ (includes Docker Scout)
- Or: Docker Engine + Docker Scout CLI plugin

## Installation

### From Source

```bash
go install github.com/ry-ops/eagle-scout/cmd/eagle-scout@latest
```

### Docker

```bash
docker pull ryops/eagle-scout:latest
```

### Binary Release

Download from [Releases](https://github.com/ry-ops/eagle-scout/releases).

## Usage

### Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "eagle-scout": {
      "command": "eagle-scout"
    }
  }
}
```

Or with Docker:

```json
{
  "mcpServers": {
    "eagle-scout": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "-v", "/var/run/docker.sock:/var/run/docker.sock", "ryops/eagle-scout"]
    }
  }
}
```

## MCP Tools

| Tool | Description |
|------|-------------|
| `scout_cves` | Scan image for CVEs with severity filtering |
| `scout_quickview` | Quick security overview of an image |
| `scout_compare` | Compare two images for security differences |
| `scout_sbom` | Generate SBOM (SPDX, CycloneDX, JSON) |
| `scout_recommendations` | Get base image update suggestions |
| `scout_policy` | Evaluate images against security policies |
| `scout_attestation` | Manage attestations on images |
| `scout_repo` | Enable/disable Scout on repositories |
| `scout_vex` | Manage VEX statements |
| `scout_version` | Get Docker Scout version info |

## Examples

### Scan an image for CVEs

```
> Use scout_cves to scan ryops/aiana:latest for critical vulnerabilities
```

### Compare image versions

```
> Use scout_compare to see what changed between ryops/aiana:v1.0.0 and ryops/aiana:latest
```

### Generate SBOM

```
> Use scout_sbom to generate a CycloneDX SBOM for my-app:latest
```

### Get update recommendations

```
> Use scout_recommendations to see if there's a better base image for my-app:latest
```

## Fabric Ecosystem

eagle-scout is part of the ry-ops fabric:

| Fabric | Language | Role |
|--------|----------|------|
| [git-steer](https://github.com/ry-ops/git-steer) | TypeScript | GitHub repo management |
| [aiana](https://github.com/ry-ops/aiana) | Python | Semantic memory |
| [n8n-fabric](https://github.com/ry-ops/n8n-fabric) | Python | Workflow automation |
| **eagle-scout** | Go | Container security |

## Development

```bash
# Clone
git clone https://github.com/ry-ops/eagle-scout
cd eagle-scout

# Build
go build -o eagle-scout ./cmd/eagle-scout

# Run
./eagle-scout

# Test
go test ./...
```

## Security

eagle-scout v1.1.1 includes important security updates:
- **Go 1.25.7** - Fixes 3 stdlib vulnerabilities (including 1 HIGH severity)
- **Updated base image** - Latest Docker CLI with security patches
- **Vulnerability reduction** - 30% reduction (20 → 14), 100% HIGH severity elimination

See [SECURITY_FIXES.md](SECURITY_FIXES.md) for details.

## Docker Hub Auto-Build

This repository is connected to Docker Hub. Any push to `main` automatically triggers a new build:
- Latest commit → `ryops/eagle-scout:latest`
- Version tags → `ryops/eagle-scout:v1.1.1`

## License

MIT License - see [LICENSE](LICENSE) file.

---

**Docker Hub:** [ryops/eagle-scout](https://hub.docker.com/r/ryops/eagle-scout)

**Version:** 1.1.1

**Status:** Alpha
