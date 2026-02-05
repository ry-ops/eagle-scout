// Package mcp provides the MCP server implementation for Docker Scout
package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ry-ops/eagle-scout/internal/scout"
)

// Server is the MCP server for Docker Scout
type Server struct {
	version string
	scout   *scout.Client
}

// NewServer creates a new MCP server
func NewServer(version string) *Server {
	return &Server{version: version}
}

// JSON-RPC types
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP types
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Capabilities struct {
	Tools *ToolsCapability `json:"tools,omitempty"`
}

type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type InitializeResult struct {
	ProtocolVersion string       `json:"protocolVersion"`
	ServerInfo      ServerInfo   `json:"serverInfo"`
	Capabilities    Capabilities `json:"capabilities"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description,omitempty"`
	Enum        []string `json:"enum,omitempty"`
	Default     any      `json:"default,omitempty"`
}

type ToolsListResult struct {
	Tools []Tool `json:"tools"`
}

type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type ToolResult struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Define all tools
var tools = []Tool{
	{
		Name:        "scout_cves",
		Description: "Scan a container image for CVEs (Common Vulnerabilities and Exposures)",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]Property{
				"image": {
					Type:        "string",
					Description: "Image to scan (e.g., 'alpine:latest', 'ryops/aiana:latest')",
				},
				"only_fixed": {
					Type:        "boolean",
					Description: "Only show CVEs that have a fix available",
					Default:     false,
				},
				"only_severity": {
					Type:        "string",
					Description: "Filter by severity level",
					Enum:        []string{"critical", "high", "medium", "low"},
				},
				"platform": {
					Type:        "string",
					Description: "Platform to scan (e.g., 'linux/amd64', 'linux/arm64')",
				},
			},
			Required: []string{"image"},
		},
	},
	{
		Name:        "scout_quickview",
		Description: "Get a quick security overview of a container image",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]Property{
				"image": {
					Type:        "string",
					Description: "Image to analyze",
				},
			},
			Required: []string{"image"},
		},
	},
	{
		Name:        "scout_compare",
		Description: "Compare security profiles of two container images",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]Property{
				"base_image": {
					Type:        "string",
					Description: "Base image to compare from (e.g., previous version)",
				},
				"target_image": {
					Type:        "string",
					Description: "Target image to compare to (e.g., new version)",
				},
				"only_fixed": {
					Type:        "boolean",
					Description: "Only show differences for fixable CVEs",
					Default:     false,
				},
			},
			Required: []string{"base_image", "target_image"},
		},
	},
	{
		Name:        "scout_sbom",
		Description: "Generate a Software Bill of Materials (SBOM) for a container image",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]Property{
				"image": {
					Type:        "string",
					Description: "Image to generate SBOM for",
				},
				"format": {
					Type:        "string",
					Description: "Output format",
					Enum:        []string{"spdx", "cyclonedx", "json"},
					Default:     "spdx",
				},
				"platform": {
					Type:        "string",
					Description: "Platform (e.g., 'linux/amd64')",
				},
			},
			Required: []string{"image"},
		},
	},
	{
		Name:        "scout_recommendations",
		Description: "Get base image update recommendations and remediation suggestions",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]Property{
				"image": {
					Type:        "string",
					Description: "Image to get recommendations for",
				},
				"only_refresh": {
					Type:        "boolean",
					Description: "Only show refresh recommendations (same tag)",
					Default:     false,
				},
				"only_update": {
					Type:        "boolean",
					Description: "Only show update recommendations (newer tag)",
					Default:     false,
				},
				"tag": {
					Type:        "string",
					Description: "Specific tag to recommend",
				},
			},
			Required: []string{"image"},
		},
	},
	{
		Name:        "scout_policy",
		Description: "Evaluate security policies against a container image",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]Property{
				"image": {
					Type:        "string",
					Description: "Image to evaluate",
				},
				"env": {
					Type:        "string",
					Description: "Environment for policy evaluation",
				},
				"org": {
					Type:        "string",
					Description: "Docker organization",
				},
			},
			Required: []string{"image"},
		},
	},
	{
		Name:        "scout_attestation",
		Description: "Manage attestations on container images",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]Property{
				"action": {
					Type:        "string",
					Description: "Action to perform",
					Enum:        []string{"add", "list"},
				},
				"image": {
					Type:        "string",
					Description: "Image to manage attestations for",
				},
				"file": {
					Type:        "string",
					Description: "Attestation file path (for add)",
				},
				"predicate_type": {
					Type:        "string",
					Description: "Predicate type for attestation",
				},
			},
			Required: []string{"action", "image"},
		},
	},
	{
		Name:        "scout_repo",
		Description: "Manage Docker Scout on repositories",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]Property{
				"action": {
					Type:        "string",
					Description: "Action to perform",
					Enum:        []string{"list", "enable", "disable"},
				},
				"repo": {
					Type:        "string",
					Description: "Repository name (for enable/disable)",
				},
				"org": {
					Type:        "string",
					Description: "Docker organization",
				},
			},
			Required: []string{"action"},
		},
	},
	{
		Name:        "scout_vex",
		Description: "Manage VEX (Vulnerability Exploitability eXchange) statements",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]Property{
				"action": {
					Type:        "string",
					Description: "Action to perform",
					Enum:        []string{"add"},
				},
				"image": {
					Type:        "string",
					Description: "Image to manage VEX for",
				},
				"cve": {
					Type:        "string",
					Description: "CVE ID (e.g., 'CVE-2024-1234')",
				},
				"status": {
					Type:        "string",
					Description: "VEX status",
					Enum:        []string{"not_affected", "affected", "fixed", "under_investigation"},
				},
				"justification": {
					Type:        "string",
					Description: "Justification for the status",
				},
				"file": {
					Type:        "string",
					Description: "VEX file to add",
				},
			},
			Required: []string{"action", "image"},
		},
	},
	{
		Name:        "scout_version",
		Description: "Get Docker Scout version information",
		InputSchema: InputSchema{
			Type:       "object",
			Properties: map[string]Property{},
		},
	},
}

// Run starts the MCP server on stdio
func (s *Server) Run() error {
	// Initialize scout client
	client, err := scout.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Docker Scout: %w", err)
	}
	s.scout = client

	scanner := bufio.NewScanner(os.Stdin)
	// Increase buffer size for large messages
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.sendError(nil, -32700, "Parse error", err.Error())
			continue
		}

		s.handleRequest(&req)
	}

	return scanner.Err()
}

func (s *Server) handleRequest(req *Request) {
	switch req.Method {
	case "initialize":
		s.handleInitialize(req)
	case "initialized":
		// Notification, no response needed
	case "tools/list":
		s.handleToolsList(req)
	case "tools/call":
		s.handleToolsCall(req)
	default:
		s.sendError(req.ID, -32601, "Method not found", req.Method)
	}
}

func (s *Server) handleInitialize(req *Request) {
	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		ServerInfo: ServerInfo{
			Name:    "eagle-scout",
			Version: s.version,
		},
		Capabilities: Capabilities{
			Tools: &ToolsCapability{},
		},
	}
	s.sendResult(req.ID, result)
}

func (s *Server) handleToolsList(req *Request) {
	s.sendResult(req.ID, ToolsListResult{Tools: tools})
}

func (s *Server) handleToolsCall(req *Request) {
	var params CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, -32602, "Invalid params", err.Error())
		return
	}

	result, err := s.executeTool(params.Name, params.Arguments)
	if err != nil {
		s.sendResult(req.ID, ToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error: %v", err)}},
			IsError: true,
		})
		return
	}

	// Marshal result to JSON string
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		s.sendError(req.ID, -32603, "Internal error", err.Error())
		return
	}

	s.sendResult(req.ID, ToolResult{
		Content: []Content{{Type: "text", Text: string(jsonBytes)}},
	})
}

func (s *Server) executeTool(name string, args map[string]interface{}) (interface{}, error) {
	switch name {
	case "scout_cves":
		image, _ := args["image"].(string)
		if image == "" {
			return nil, fmt.Errorf("image is required")
		}
		opts := scout.CVEsOptions{
			OnlyFixed:    getBool(args, "only_fixed"),
			OnlySeverity: getString(args, "only_severity"),
			Platform:     getString(args, "platform"),
		}
		return s.scout.CVEs(image, opts)

	case "scout_quickview":
		image, _ := args["image"].(string)
		if image == "" {
			return nil, fmt.Errorf("image is required")
		}
		return s.scout.Quickview(image)

	case "scout_compare":
		baseImage, _ := args["base_image"].(string)
		targetImage, _ := args["target_image"].(string)
		if baseImage == "" || targetImage == "" {
			return nil, fmt.Errorf("base_image and target_image are required")
		}
		opts := scout.CompareOptions{
			OnlyFixed: getBool(args, "only_fixed"),
		}
		return s.scout.Compare(baseImage, targetImage, opts)

	case "scout_sbom":
		image, _ := args["image"].(string)
		if image == "" {
			return nil, fmt.Errorf("image is required")
		}
		opts := scout.SBOMOptions{
			Format:   getString(args, "format"),
			Platform: getString(args, "platform"),
		}
		return s.scout.SBOM(image, opts)

	case "scout_recommendations":
		image, _ := args["image"].(string)
		if image == "" {
			return nil, fmt.Errorf("image is required")
		}
		opts := scout.RecommendationsOptions{
			OnlyRefresh: getBool(args, "only_refresh"),
			OnlyUpdate:  getBool(args, "only_update"),
			Tag:         getString(args, "tag"),
		}
		return s.scout.Recommendations(image, opts)

	case "scout_policy":
		image, _ := args["image"].(string)
		if image == "" {
			return nil, fmt.Errorf("image is required")
		}
		opts := scout.PolicyOptions{
			Env: getString(args, "env"),
			Org: getString(args, "org"),
		}
		return s.scout.Policy(image, opts)

	case "scout_attestation":
		action, _ := args["action"].(string)
		image, _ := args["image"].(string)
		if action == "" || image == "" {
			return nil, fmt.Errorf("action and image are required")
		}
		opts := scout.AttestationOptions{
			File:          getString(args, "file"),
			PredicateType: getString(args, "predicate_type"),
		}
		switch action {
		case "add":
			return s.scout.AttestationAdd(image, opts)
		default:
			return nil, fmt.Errorf("unsupported attestation action: %s", action)
		}

	case "scout_repo":
		action, _ := args["action"].(string)
		if action == "" {
			return nil, fmt.Errorf("action is required")
		}
		opts := scout.RepoOptions{
			Org: getString(args, "org"),
		}
		switch action {
		case "list":
			return s.scout.RepoList(opts)
		case "enable":
			repo, _ := args["repo"].(string)
			if repo == "" {
				return nil, fmt.Errorf("repo is required for enable")
			}
			return s.scout.RepoEnable(repo, opts)
		case "disable":
			repo, _ := args["repo"].(string)
			if repo == "" {
				return nil, fmt.Errorf("repo is required for disable")
			}
			return s.scout.RepoDisable(repo, opts)
		default:
			return nil, fmt.Errorf("unsupported repo action: %s", action)
		}

	case "scout_vex":
		action, _ := args["action"].(string)
		image, _ := args["image"].(string)
		if action == "" || image == "" {
			return nil, fmt.Errorf("action and image are required")
		}
		opts := scout.VexOptions{
			File:          getString(args, "file"),
			CVE:           getString(args, "cve"),
			Status:        getString(args, "status"),
			Justification: getString(args, "justification"),
		}
		switch action {
		case "add":
			return s.scout.VexAdd(image, opts)
		default:
			return nil, fmt.Errorf("unsupported vex action: %s", action)
		}

	case "scout_version":
		return s.scout.Version()

	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

func (s *Server) sendResult(id interface{}, result interface{}) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	s.send(resp)
}

func (s *Server) sendError(id interface{}, code int, message string, data interface{}) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &Error{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	s.send(resp)
}

func (s *Server) send(resp Response) {
	bytes, _ := json.Marshal(resp)
	fmt.Println(string(bytes))
}

// Helper functions
func getString(args map[string]interface{}, key string) string {
	if v, ok := args[key].(string); ok {
		return v
	}
	return ""
}

func getBool(args map[string]interface{}, key string) bool {
	if v, ok := args[key].(bool); ok {
		return v
	}
	return false
}
