// eagle-scout - MCP server for Docker Scout
package main

import (
	"fmt"
	"os"

	"github.com/ry-ops/eagle-scout/internal/mcp"
)

var version = "1.1.0"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "--version", "-v":
			fmt.Printf("eagle-scout version %s\n", version)
			return
		case "help", "--help", "-h":
			printHelp()
			return
		}
	}

	// Start MCP server on stdio
	server := mcp.NewServer(version)
	if err := server.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`eagle-scout - Docker Scout MCP Server

Usage:
  eagle-scout              Start MCP server (stdio transport)
  eagle-scout version      Show version
  eagle-scout help         Show this help

Environment:
  DOCKER_SCOUT_HUB_USER  Docker Hub username (optional)
  DOCKER_SCOUT_HUB_PAT   Docker Hub PAT (optional)

The MCP server exposes Docker Scout functionality as tools:
  - scout_cves           Scan image for CVEs
  - scout_quickview      Quick security overview
  - scout_compare        Compare two images
  - scout_sbom           Generate SBOM
  - scout_recommendations Get remediation suggestions
  - scout_policy         Evaluate security policies
  - scout_attestation    Manage attestations
  - scout_repo           Manage Scout-enabled repos
  - scout_vex            Manage VEX statements (add/list)
  - scout_environment    Manage environments (list/set)
  - scout_cache          Manage local cache (df/prune)
  - scout_enroll         Enroll organization
  - scout_watch          Enable/disable continuous monitoring

For more info: https://github.com/ry-ops/eagle-scout`)
}
