# scout-mcp

**MCP Server for Docker Scout** - Container security scanning via Model Context Protocol.

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
go install github.com/ry-ops/scout-mcp/cmd/scout-mcp@latest
```

### Docker

```bash
docker pull ryops/scout-mcp:latest
```

### Binary Release

Download from [Releases](https://github.com/ry-ops/scout-mcp/releases).

## Usage

### Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "scout-mcp": {
      "command": "scout-mcp"
    }
  }
}
```

Or with Docker:

```json
{
  "mcpServers": {
    "scout-mcp": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "-v", "/var/run/docker.sock:/var/run/docker.sock", "ryops/scout-mcp"]
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

scout-mcp is part of the ry-ops fabric:

| Fabric | Language | Role |
|--------|----------|------|
| [git-steer](https://github.com/ry-ops/git-steer) | TypeScript | GitHub repo management |
| [aiana](https://github.com/ry-ops/aiana) | Python | Semantic memory |
| [n8n-fabric](https://github.com/ry-ops/n8n-fabric) | Python | Workflow automation |
| **scout-mcp** | Go | Container security |

## Development

```bash
# Clone
git clone https://github.com/ry-ops/scout-mcp
cd scout-mcp

# Build
go build -o scout-mcp ./cmd/scout-mcp

# Run
./scout-mcp

# Test
go test ./...
```

## License

MIT License - see [LICENSE](LICENSE) file.

---

**Docker Hub:** [ryops/scout-mcp](https://hub.docker.com/r/ryops/scout-mcp)

**Status:** Alpha
