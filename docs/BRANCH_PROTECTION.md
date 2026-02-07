# Branch Protection Setup

This document describes the required branch protection rules for the `main` branch.

## GitHub Settings

Navigate to: **Settings → Branches → Add branch protection rule**

### Branch name pattern
```
main
```

### Required Settings

#### Protect matching branches

- [x] **Require a pull request before merging**
  - [x] Require approvals: `1`
  - [x] Dismiss stale pull request approvals when new commits are pushed
  - [x] Require approval of the most recent reviewable push

- [x] **Require status checks to pass before merging**
  - [x] Require branches to be up to date before merging
  - Required status checks:
    - `Build & Test`
    - `Security Scan`
    - `Policy Check`

- [x] **Require conversation resolution before merging**

- [x] **Require signed commits** (recommended)

- [x] **Require linear history**
  - Prevents merge commits, enforces squash or rebase

- [x] **Include administrators**
  - Enforces rules on admins — prevents bypassing security gates

#### Rules applied to everyone including administrators

- [x] **Do not allow bypassing the above settings**

- [x] **Restrict who can push to matching branches**
  - Only allow merges through PR

- [x] **Block force pushes**

- [x] **Block deletions**

## CLI Setup (using gh)

```bash
# Set branch protection via GitHub CLI
gh api repos/ry-ops/eagle-scout/branches/main/protection \
  --method PUT \
  --field required_status_checks='{"strict":true,"contexts":["Build & Test","Security Scan","Policy Check"]}' \
  --field enforce_admins=true \
  --field required_pull_request_reviews='{"dismiss_stale_reviews":true,"require_code_owner_reviews":false,"required_approving_review_count":1}' \
  --field restrictions=null \
  --field required_linear_history=true \
  --field allow_force_pushes=false \
  --field allow_deletions=false
```

## Workflow

```
Developer                    GitHub                         Main Branch
    |                           |                               |
    |-- Create feature branch --|                               |
    |-- Push commits ---------->|                               |
    |-- Open PR --------------->|                               |
    |                           |-- Run CI ------------------>  |
    |                           |   - Build & Test              |
    |                           |   - Security Scan             |
    |                           |   - Policy Check              |
    |                           |                               |
    |                           |<-- All checks pass            |
    |                           |                               |
    |<-- Request review --------|                               |
    |                           |                               |
    |-- Address feedback ------>|                               |
    |                           |-- Re-run CI                   |
    |                           |                               |
    |                           |<-- Approved + Checks pass     |
    |                           |                               |
    |                           |-- Squash & Merge ------------>|
    |                           |                               |
    |                           |-- Trigger publish workflow -->|
    |                           |                               |
```

## Required Secrets

Configure in **Settings → Secrets and variables → Actions**:

| Secret | Description |
|--------|-------------|
| `DOCKERHUB_USERNAME` | Docker Hub username |
| `DOCKERHUB_TOKEN` | Docker Hub access token |

## Rulesets (Alternative)

GitHub Rulesets provide more granular control. Navigate to **Settings → Rules → Rulesets** to create:

```yaml
name: main-protection
target: branch
enforcement: active
conditions:
  ref_name:
    include: ["refs/heads/main"]
rules:
  - type: pull_request
    parameters:
      required_approving_review_count: 1
      dismiss_stale_reviews_on_push: true
      require_last_push_approval: true
  - type: required_status_checks
    parameters:
      strict_required_status_checks_policy: true
      required_status_checks:
        - context: "Build & Test"
        - context: "Security Scan"
        - context: "Policy Check"
  - type: non_fast_forward
  - type: deletion
```
