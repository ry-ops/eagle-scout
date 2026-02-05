// Package scout provides a client wrapper for Docker Scout CLI
package scout

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Client wraps Docker Scout CLI operations
type Client struct {
	dockerPath string
}

// NewClient creates a new Scout client
func NewClient() (*Client, error) {
	// Find docker executable
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		return nil, fmt.Errorf("docker not found in PATH: %w", err)
	}

	// Verify scout is available
	cmd := exec.Command(dockerPath, "scout", "version")
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("docker scout not available: %w", err)
	}

	return &Client{dockerPath: dockerPath}, nil
}

// CVEResult represents a CVE finding
type CVEResult struct {
	ID          string `json:"id"`
	Severity    string `json:"severity"`
	CVSS        float64 `json:"cvss,omitempty"`
	Package     string `json:"package"`
	Version     string `json:"version"`
	FixVersion  string `json:"fix_version,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

// CVEsOutput represents the output of scout cves
type CVEsOutput struct {
	Image          string      `json:"image"`
	TotalVulns     int         `json:"total_vulnerabilities"`
	Critical       int         `json:"critical"`
	High           int         `json:"high"`
	Medium         int         `json:"medium"`
	Low            int         `json:"low"`
	Unspecified    int         `json:"unspecified"`
	Vulnerabilities []CVEResult `json:"vulnerabilities,omitempty"`
	RawOutput      string      `json:"raw_output,omitempty"`
}

// CVEs scans an image for CVEs
func (c *Client) CVEs(image string, opts CVEsOptions) (*CVEsOutput, error) {
	args := []string{"scout", "cves", image}

	if opts.OnlyFixed {
		args = append(args, "--only-fixed")
	}
	if opts.OnlySeverity != "" {
		args = append(args, "--only-severity", opts.OnlySeverity)
	}
	if opts.Format != "" {
		args = append(args, "--format", opts.Format)
	}
	if opts.ExitCode {
		args = append(args, "--exit-code")
	}
	if opts.Platform != "" {
		args = append(args, "--platform", opts.Platform)
	}

	output, err := c.run(args...)
	if err != nil {
		// Scout returns non-zero if vulnerabilities found with --exit-code
		if opts.ExitCode && output != "" {
			// Parse the output anyway
		} else {
			return nil, err
		}
	}

	result := &CVEsOutput{
		Image:     image,
		RawOutput: output,
	}

	// Parse summary counts from output
	c.parseCVECounts(output, result)

	return result, nil
}

// CVEsOptions for CVEs command
type CVEsOptions struct {
	OnlyFixed    bool
	OnlySeverity string // critical, high, medium, low
	Format       string // json, sarif, spdx, markdown
	ExitCode     bool   // Return non-zero if vulns found
	Platform     string // linux/amd64, linux/arm64
}

// QuickviewOutput represents scout quickview output
type QuickviewOutput struct {
	Image       string `json:"image"`
	BaseImage   string `json:"base_image,omitempty"`
	Packages    int    `json:"packages"`
	Layers      int    `json:"layers"`
	TotalVulns  int    `json:"total_vulnerabilities"`
	Critical    int    `json:"critical"`
	High        int    `json:"high"`
	Medium      int    `json:"medium"`
	Low         int    `json:"low"`
	RawOutput   string `json:"raw_output"`
}

// Quickview provides a quick security overview
func (c *Client) Quickview(image string) (*QuickviewOutput, error) {
	output, err := c.run("scout", "quickview", image)
	if err != nil {
		return nil, err
	}

	result := &QuickviewOutput{
		Image:     image,
		RawOutput: output,
	}

	// Parse the output
	c.parseQuickview(output, result)

	return result, nil
}

// CompareOutput represents scout compare output
type CompareOutput struct {
	BaseImage    string `json:"base_image"`
	TargetImage  string `json:"target_image"`
	AddedVulns   int    `json:"added_vulnerabilities"`
	RemovedVulns int    `json:"removed_vulnerabilities"`
	RawOutput    string `json:"raw_output"`
}

// Compare compares two images
func (c *Client) Compare(baseImage, targetImage string, opts CompareOptions) (*CompareOutput, error) {
	args := []string{"scout", "compare", "--to", baseImage, targetImage}

	if opts.OnlyFixed {
		args = append(args, "--only-fixed")
	}
	if opts.ExitCode {
		args = append(args, "--exit-code")
	}

	output, err := c.run(args...)
	if err != nil && output == "" {
		return nil, err
	}

	return &CompareOutput{
		BaseImage:   baseImage,
		TargetImage: targetImage,
		RawOutput:   output,
	}, nil
}

// CompareOptions for Compare command
type CompareOptions struct {
	OnlyFixed bool
	ExitCode  bool
}

// SBOMOutput represents scout sbom output
type SBOMOutput struct {
	Image     string `json:"image"`
	Format    string `json:"format"`
	RawOutput string `json:"raw_output"`
}

// SBOM generates a Software Bill of Materials
func (c *Client) SBOM(image string, opts SBOMOptions) (*SBOMOutput, error) {
	args := []string{"scout", "sbom", image}

	format := opts.Format
	if format == "" {
		format = "spdx"
	}
	args = append(args, "--format", format)

	if opts.Platform != "" {
		args = append(args, "--platform", opts.Platform)
	}

	output, err := c.run(args...)
	if err != nil {
		return nil, err
	}

	return &SBOMOutput{
		Image:     image,
		Format:    format,
		RawOutput: output,
	}, nil
}

// SBOMOptions for SBOM command
type SBOMOptions struct {
	Format   string // spdx, cyclonedx, json
	Platform string
}

// RecommendationsOutput represents scout recommendations output
type RecommendationsOutput struct {
	Image              string   `json:"image"`
	CurrentBase        string   `json:"current_base,omitempty"`
	RecommendedBase    string   `json:"recommended_base,omitempty"`
	VulnsReduced       int      `json:"vulnerabilities_reduced,omitempty"`
	Recommendations    []string `json:"recommendations,omitempty"`
	RawOutput          string   `json:"raw_output"`
}

// Recommendations gets base image update recommendations
func (c *Client) Recommendations(image string, opts RecommendationsOptions) (*RecommendationsOutput, error) {
	args := []string{"scout", "recommendations", image}

	if opts.OnlyRefresh {
		args = append(args, "--only-refresh")
	}
	if opts.OnlyUpdate {
		args = append(args, "--only-update")
	}
	if opts.Tag != "" {
		args = append(args, "--tag", opts.Tag)
	}

	output, err := c.run(args...)
	if err != nil {
		return nil, err
	}

	return &RecommendationsOutput{
		Image:     image,
		RawOutput: output,
	}, nil
}

// RecommendationsOptions for Recommendations command
type RecommendationsOptions struct {
	OnlyRefresh bool
	OnlyUpdate  bool
	Tag         string
}

// PolicyOutput represents scout policy output
type PolicyOutput struct {
	Image      string          `json:"image"`
	Passed     bool            `json:"passed"`
	Results    []PolicyResult  `json:"results,omitempty"`
	RawOutput  string          `json:"raw_output"`
}

// PolicyResult represents a single policy evaluation result
type PolicyResult struct {
	Policy  string `json:"policy"`
	Status  string `json:"status"` // passed, failed, warning
	Message string `json:"message,omitempty"`
}

// Policy evaluates policies against an image
func (c *Client) Policy(image string, opts PolicyOptions) (*PolicyOutput, error) {
	args := []string{"scout", "policy", image}

	if opts.Env != "" {
		args = append(args, "--env", opts.Env)
	}
	if opts.Org != "" {
		args = append(args, "--org", opts.Org)
	}
	if opts.ExitCode {
		args = append(args, "--exit-code")
	}

	output, err := c.run(args...)
	if err != nil && output == "" {
		return nil, err
	}

	result := &PolicyOutput{
		Image:     image,
		RawOutput: output,
		Passed:    err == nil,
	}

	return result, nil
}

// PolicyOptions for Policy command
type PolicyOptions struct {
	Env      string
	Org      string
	ExitCode bool
}

// AttestationOutput represents attestation operations output
type AttestationOutput struct {
	Image     string   `json:"image"`
	Action    string   `json:"action"`
	Success   bool     `json:"success"`
	Details   []string `json:"details,omitempty"`
	RawOutput string   `json:"raw_output"`
}

// AttestationAdd adds an attestation to an image
func (c *Client) AttestationAdd(image string, opts AttestationOptions) (*AttestationOutput, error) {
	args := []string{"scout", "attestation", "add", image}

	if opts.File != "" {
		args = append(args, "--file", opts.File)
	}
	if opts.PredicateType != "" {
		args = append(args, "--predicate-type", opts.PredicateType)
	}

	output, err := c.run(args...)
	if err != nil {
		return nil, err
	}

	return &AttestationOutput{
		Image:     image,
		Action:    "add",
		Success:   true,
		RawOutput: output,
	}, nil
}

// AttestationOptions for Attestation commands
type AttestationOptions struct {
	File          string
	PredicateType string
}

// RepoOutput represents scout repo operations output
type RepoOutput struct {
	Action    string   `json:"action"`
	Repos     []string `json:"repos,omitempty"`
	Success   bool     `json:"success"`
	RawOutput string   `json:"raw_output"`
}

// RepoList lists Scout-enabled repositories
func (c *Client) RepoList(opts RepoOptions) (*RepoOutput, error) {
	args := []string{"scout", "repo", "list"}

	if opts.Org != "" {
		args = append(args, "--org", opts.Org)
	}

	output, err := c.run(args...)
	if err != nil {
		return nil, err
	}

	return &RepoOutput{
		Action:    "list",
		Success:   true,
		RawOutput: output,
	}, nil
}

// RepoEnable enables Scout on a repository
func (c *Client) RepoEnable(repo string, opts RepoOptions) (*RepoOutput, error) {
	args := []string{"scout", "repo", "enable", repo}

	if opts.Org != "" {
		args = append(args, "--org", opts.Org)
	}

	output, err := c.run(args...)
	if err != nil {
		return nil, err
	}

	return &RepoOutput{
		Action:    "enable",
		Repos:     []string{repo},
		Success:   true,
		RawOutput: output,
	}, nil
}

// RepoDisable disables Scout on a repository
func (c *Client) RepoDisable(repo string, opts RepoOptions) (*RepoOutput, error) {
	args := []string{"scout", "repo", "disable", repo}

	if opts.Org != "" {
		args = append(args, "--org", opts.Org)
	}

	output, err := c.run(args...)
	if err != nil {
		return nil, err
	}

	return &RepoOutput{
		Action:    "disable",
		Repos:     []string{repo},
		Success:   true,
		RawOutput: output,
	}, nil
}

// RepoOptions for Repo commands
type RepoOptions struct {
	Org string
}

// VexOutput represents VEX operations output
type VexOutput struct {
	Image     string `json:"image"`
	Action    string `json:"action"`
	Success   bool   `json:"success"`
	RawOutput string `json:"raw_output"`
}

// VexAdd adds a VEX statement
func (c *Client) VexAdd(image string, opts VexOptions) (*VexOutput, error) {
	args := []string{"scout", "vex", "add", image}

	if opts.File != "" {
		args = append(args, "--file", opts.File)
	}
	if opts.CVE != "" {
		args = append(args, "--cve", opts.CVE)
	}
	if opts.Status != "" {
		args = append(args, "--status", opts.Status)
	}
	if opts.Justification != "" {
		args = append(args, "--justification", opts.Justification)
	}

	output, err := c.run(args...)
	if err != nil {
		return nil, err
	}

	return &VexOutput{
		Image:     image,
		Action:    "add",
		Success:   true,
		RawOutput: output,
	}, nil
}

// VexList lists VEX statements for an image
func (c *Client) VexList(image string, opts VexOptions) (*VexOutput, error) {
	args := []string{"scout", "cves", image, "--vex-location", "image"}

	output, err := c.run(args...)
	if err != nil && output == "" {
		return nil, err
	}

	return &VexOutput{
		Image:     image,
		Action:    "list",
		Success:   true,
		RawOutput: output,
	}, nil
}

// VexOptions for VEX commands
type VexOptions struct {
	File          string
	CVE           string
	Status        string // not_affected, affected, fixed, under_investigation
	Justification string
}

// EnvironmentOutput represents environment operations output
type EnvironmentOutput struct {
	Action       string   `json:"action"`
	Environment  string   `json:"environment,omitempty"`
	Image        string   `json:"image,omitempty"`
	Environments []string `json:"environments,omitempty"`
	Success      bool     `json:"success"`
	RawOutput    string   `json:"raw_output"`
}

// EnvironmentList lists environments
func (c *Client) EnvironmentList(opts EnvironmentOptions) (*EnvironmentOutput, error) {
	args := []string{"scout", "environment"}

	if opts.Org != "" {
		args = append(args, "--org", opts.Org)
	}

	output, err := c.run(args...)
	if err != nil {
		return nil, err
	}

	return &EnvironmentOutput{
		Action:    "list",
		Success:   true,
		RawOutput: output,
	}, nil
}

// EnvironmentSet sets an image in an environment
func (c *Client) EnvironmentSet(env, image string, opts EnvironmentOptions) (*EnvironmentOutput, error) {
	args := []string{"scout", "environment", env, image}

	if opts.Org != "" {
		args = append(args, "--org", opts.Org)
	}

	output, err := c.run(args...)
	if err != nil {
		return nil, err
	}

	return &EnvironmentOutput{
		Action:      "set",
		Environment: env,
		Image:       image,
		Success:     true,
		RawOutput:   output,
	}, nil
}

// EnvironmentOptions for Environment commands
type EnvironmentOptions struct {
	Org string
}

// CacheOutput represents cache operations output
type CacheOutput struct {
	Action    string `json:"action"`
	Success   bool   `json:"success"`
	Size      string `json:"size,omitempty"`
	RawOutput string `json:"raw_output"`
}

// CachePrune prunes the local cache
func (c *Client) CachePrune() (*CacheOutput, error) {
	output, err := c.run("scout", "cache", "prune", "--force")
	if err != nil {
		return nil, err
	}

	return &CacheOutput{
		Action:    "prune",
		Success:   true,
		RawOutput: output,
	}, nil
}

// CacheDF shows cache disk usage
func (c *Client) CacheDF() (*CacheOutput, error) {
	output, err := c.run("scout", "cache", "df")
	if err != nil {
		return nil, err
	}

	return &CacheOutput{
		Action:    "df",
		Success:   true,
		RawOutput: output,
	}, nil
}

// EnrollOutput represents enrollment output
type EnrollOutput struct {
	Org       string `json:"org"`
	Success   bool   `json:"success"`
	RawOutput string `json:"raw_output"`
}

// Enroll enrolls an organization with Docker Scout
func (c *Client) Enroll(org string) (*EnrollOutput, error) {
	args := []string{"scout", "enroll", org}

	output, err := c.run(args...)
	if err != nil {
		return nil, err
	}

	return &EnrollOutput{
		Org:       org,
		Success:   true,
		RawOutput: output,
	}, nil
}

// WatchOutput represents watch configuration output
type WatchOutput struct {
	Repository string `json:"repository"`
	Action     string `json:"action"`
	Success    bool   `json:"success"`
	RawOutput  string `json:"raw_output"`
}

// WatchEnable enables continuous monitoring for a repository
func (c *Client) WatchEnable(repo string, opts WatchOptions) (*WatchOutput, error) {
	// Watch is enabled via repo enable with integration
	args := []string{"scout", "repo", "enable", repo}

	if opts.Org != "" {
		args = append(args, "--org", opts.Org)
	}
	if opts.Integration != "" {
		args = append(args, "--integration", opts.Integration)
	}

	output, err := c.run(args...)
	if err != nil {
		return nil, err
	}

	return &WatchOutput{
		Repository: repo,
		Action:     "enable",
		Success:    true,
		RawOutput:  output,
	}, nil
}

// WatchDisable disables continuous monitoring
func (c *Client) WatchDisable(repo string, opts WatchOptions) (*WatchOutput, error) {
	args := []string{"scout", "repo", "disable", repo}

	if opts.Org != "" {
		args = append(args, "--org", opts.Org)
	}

	output, err := c.run(args...)
	if err != nil {
		return nil, err
	}

	return &WatchOutput{
		Repository: repo,
		Action:     "disable",
		Success:    true,
		RawOutput:  output,
	}, nil
}

// WatchOptions for Watch commands
type WatchOptions struct {
	Org         string
	Integration string // github, gitlab, etc.
}

// VersionInfo represents scout version information
type VersionInfo struct {
	Version   string `json:"version"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
}

// Version returns Docker Scout version info
func (c *Client) Version() (*VersionInfo, error) {
	output, err := c.run("scout", "version")
	if err != nil {
		return nil, err
	}

	return &VersionInfo{
		Version: strings.TrimSpace(output),
	}, nil
}

// run executes a docker command and returns output
func (c *Client) run(args ...string) (string, error) {
	cmd := exec.Command(c.dockerPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String()

	if err != nil {
		if stderr.Len() > 0 {
			return output, fmt.Errorf("%s: %s", err, stderr.String())
		}
		return output, err
	}

	return output, nil
}

// runJSON executes a command and parses JSON output
func (c *Client) runJSON(v interface{}, args ...string) error {
	output, err := c.run(args...)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(output), v)
}

// parseCVECounts parses CVE counts from output
func (c *Client) parseCVECounts(output string, result *CVEsOutput) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.ToLower(line)
		if strings.Contains(line, "critical") {
			fmt.Sscanf(line, "%d critical", &result.Critical)
		}
		if strings.Contains(line, "high") {
			fmt.Sscanf(line, "%d high", &result.High)
		}
		if strings.Contains(line, "medium") {
			fmt.Sscanf(line, "%d medium", &result.Medium)
		}
		if strings.Contains(line, "low") {
			fmt.Sscanf(line, "%d low", &result.Low)
		}
	}
	result.TotalVulns = result.Critical + result.High + result.Medium + result.Low
}

// parseQuickview parses quickview output
func (c *Client) parseQuickview(output string, result *QuickviewOutput) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.ToLower(line)
		if strings.Contains(line, "packages") {
			fmt.Sscanf(line, "%d packages", &result.Packages)
		}
		if strings.Contains(line, "layers") {
			fmt.Sscanf(line, "%d layers", &result.Layers)
		}
	}
	c.parseCVECounts(output, &CVEsOutput{})
}
