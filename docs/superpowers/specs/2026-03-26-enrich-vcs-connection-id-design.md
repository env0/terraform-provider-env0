# Enrich github_installation_id / bitbucketClientKey with VCS Connection ID

**Linear ticket:** ENG-1398
**Date:** 2026-03-26

## Problem

The move to VCS connections for GitHub broke backwards compatibility. Users who create a GitHub connection via the UI (which is now a VCS connection) and then use the Terraform provider with `github_installation_id` fail on deploy/template creation because the backend needs the `vcs_connection_id` for the new authorization flow.

Currently, enrichment only works one way: the backend auto-populates VCS fields (like `github_installation_id`) when the user provides `vcs_connection_id`. The reverse does not happen.

## Solution

Provider-side enrichment: before sending create/update payloads to the backend API, the provider looks up the matching VCS connection by `github_installation_id` or `bitbucket_client_key` and includes the `vcs_connection_id` in the payload.

## Design

### 1. Update VcsConnection struct

File: `client/vcs_connection.go`

Add `GithubInstallationId` and `BitbucketClientKey` fields to the `VcsConnection` struct to match the API response:

```go
type VcsConnection struct {
    Id                   string `json:"id"`
    Name                 string `json:"name"`
    Type                 string `json:"type"`
    Url                  string `json:"url"`
    VcsAgentKey          string `json:"vcsAgentKey"`
    GithubInstallationId int    `json:"githubInstallationId"`
    BitbucketClientKey   string `json:"bitbucketClientKey"`
}
```

### 2. Add enrichment helper function

File: `env0/utils.go`

```go
func enrichVcsConnectionId(apiClient client.ApiClientInterface, githubInstallationId int, bitbucketClientKey string, vcsConnectionId *string) error
```

Logic:
- If `vcsConnectionId` is already set, return immediately (no-op)
- If neither `githubInstallationId` nor `bitbucketClientKey` is set, return immediately
- Call `apiClient.VcsConnections()` to list all org VCS connections
- Filter by matching `githubInstallationId` (if non-zero) or `bitbucketClientKey` (if non-empty)
- Set `*vcsConnectionId` to the matching connection's ID
- Return error if no matching connection found

### 3. Integrate enrichment into all resource flows

Call `enrichVcsConnectionId` after `Invalidate()` but before the API call in these resources:

| Resource | Files | Operations |
|----------|-------|------------|
| Template | `env0/resource_template.go` | create, update |
| Environment (without_template_settings) | `env0/resource_environment.go` | create, update |
| Custom Flow | `env0/resource_custom_flow.go` | create, update |
| Approval Policy | `env0/resource_approval_policy.go` | create, update |
| Module | `env0/resource_module.go` | create, update |
| Environment Discovery | `env0/resource_environment_discovery.go` | put (create/update) |

For modules where `GithubInstallationId` is `*int`, dereference before calling the helper.

### 4. Tests

- Unit test the enrichment helper function in `env0/utils_test.go`
- Update existing resource tests to mock `VcsConnections()` where enrichment is triggered
- Test cases: enrichment succeeds, no match found (error), vcs_connection_id already set (no-op), neither field set (no-op)

## Key Design Decisions

- **Enrichment happens after Invalidate()**: The mutual exclusivity validation (`github_installation_id` and `vcs_connection_id` are mutually exclusive) checks user input. The provider enriches the payload afterward, so both fields are sent to the API.
- **Error on no match**: If a user provides a `github_installation_id` that has no corresponding VCS connection, the provider fails with a clear error rather than silently proceeding.
- **All resources enriched**: Every resource with VCS fields gets enrichment, not just templates and environments.
