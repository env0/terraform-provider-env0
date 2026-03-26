# Enrich VCS Connection ID Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** When a user provides `github_installation_id` or `bitbucket_client_key`, the provider enriches the payload with the matching `vcs_connection_id` before sending to the API.

**Architecture:** Add two new fields to `VcsConnection` struct, create an `enrichVcsConnectionId` helper in `utils.go`, and call it in all resource create/update flows after validation but before the API call.

**Tech Stack:** Go, Terraform Plugin SDK v2, gomock

---

### Task 1: Update VcsConnection Struct

**Files:**
- Modify: `client/vcs_connection.go:3-9`

- [ ] **Step 1: Add GithubInstallationId and BitbucketClientKey fields**

In `client/vcs_connection.go`, add two fields to the `VcsConnection` struct:

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

- [ ] **Step 2: Verify it compiles**

Run: `go build ./...`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add client/vcs_connection.go
git commit -m "feat: add GithubInstallationId and BitbucketClientKey to VcsConnection struct"
```

---

### Task 2: Add enrichVcsConnectionId Helper + Unit Tests

**Files:**
- Modify: `env0/utils.go` (add function after `suppressVcsFieldDrift` around line 539)
- Modify: `env0/utils_test.go` (add tests)

- [ ] **Step 1: Write failing tests for the enrichment helper**

Add the following tests to `env0/utils_test.go`:

```go
func TestEnrichVcsConnectionId(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := client.NewMockApiClientInterface(ctrl)

	vcsConnections := []client.VcsConnection{
		{
			Id:                   "conn-github-123",
			Name:                 "github-connection",
			GithubInstallationId: 12345,
		},
		{
			Id:                 "conn-bitbucket-456",
			Name:               "bitbucket-connection",
			BitbucketClientKey: "bb-key-abc",
		},
	}

	t.Run("enriches vcs_connection_id from github_installation_id", func(t *testing.T) {
		mock.EXPECT().VcsConnections().Times(1).Return(vcsConnections, nil)

		vcsConnectionId := ""
		err := enrichVcsConnectionId(mock, 12345, "", &vcsConnectionId)

		assert.NoError(t, err)
		assert.Equal(t, "conn-github-123", vcsConnectionId)
	})

	t.Run("enriches vcs_connection_id from bitbucket_client_key", func(t *testing.T) {
		mock.EXPECT().VcsConnections().Times(1).Return(vcsConnections, nil)

		vcsConnectionId := ""
		err := enrichVcsConnectionId(mock, 0, "bb-key-abc", &vcsConnectionId)

		assert.NoError(t, err)
		assert.Equal(t, "conn-bitbucket-456", vcsConnectionId)
	})

	t.Run("no-op when vcs_connection_id is already set", func(t *testing.T) {
		vcsConnectionId := "already-set"
		err := enrichVcsConnectionId(mock, 12345, "", &vcsConnectionId)

		assert.NoError(t, err)
		assert.Equal(t, "already-set", vcsConnectionId)
	})

	t.Run("no-op when neither github_installation_id nor bitbucket_client_key is set", func(t *testing.T) {
		vcsConnectionId := ""
		err := enrichVcsConnectionId(mock, 0, "", &vcsConnectionId)

		assert.NoError(t, err)
		assert.Equal(t, "", vcsConnectionId)
	})

	t.Run("error when no matching vcs connection found for github_installation_id", func(t *testing.T) {
		mock.EXPECT().VcsConnections().Times(1).Return(vcsConnections, nil)

		vcsConnectionId := ""
		err := enrichVcsConnectionId(mock, 99999, "", &vcsConnectionId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not find a VCS connection")
	})

	t.Run("error when no matching vcs connection found for bitbucket_client_key", func(t *testing.T) {
		mock.EXPECT().VcsConnections().Times(1).Return(vcsConnections, nil)

		vcsConnectionId := ""
		err := enrichVcsConnectionId(mock, 0, "nonexistent-key", &vcsConnectionId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not find a VCS connection")
	})

	t.Run("error when VcsConnections API call fails", func(t *testing.T) {
		mock.EXPECT().VcsConnections().Times(1).Return(nil, errors.New("api error"))

		vcsConnectionId := ""
		err := enrichVcsConnectionId(mock, 12345, "", &vcsConnectionId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "api error")
	})
}
```

Also add these imports at the top of `env0/utils_test.go` (merge with existing):

```go
import (
	"errors"

	"github.com/env0/terraform-provider-env0/client"
	"go.uber.org/mock/gomock"
)
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./env0/ -run TestEnrichVcsConnectionId -v`
Expected: FAIL — `enrichVcsConnectionId` is undefined

- [ ] **Step 3: Implement enrichVcsConnectionId**

Add the following function to `env0/utils.go` (after `suppressVcsFieldDrift`, around line 539):

```go
// enrichVcsConnectionId looks up the VCS connection ID from github_installation_id or
// bitbucket_client_key when the user hasn't explicitly set vcs_connection_id.
// This maintains backwards compatibility after the move to VCS connections for GitHub.
func enrichVcsConnectionId(apiClient client.ApiClientInterface, githubInstallationId int, bitbucketClientKey string, vcsConnectionId *string) error {
	if *vcsConnectionId != "" {
		return nil
	}

	if githubInstallationId == 0 && bitbucketClientKey == "" {
		return nil
	}

	connections, err := apiClient.VcsConnections()
	if err != nil {
		return fmt.Errorf("failed to fetch VCS connections: %w", err)
	}

	for _, conn := range connections {
		if githubInstallationId != 0 && conn.GithubInstallationId == githubInstallationId {
			*vcsConnectionId = conn.Id
			return nil
		}

		if bitbucketClientKey != "" && conn.BitbucketClientKey == bitbucketClientKey {
			*vcsConnectionId = conn.Id
			return nil
		}
	}

	return fmt.Errorf("could not find a VCS connection matching the provided github_installation_id or bitbucket_client_key")
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./env0/ -run TestEnrichVcsConnectionId -v`
Expected: All 7 tests PASS

- [ ] **Step 5: Commit**

```bash
git add env0/utils.go env0/utils_test.go
git commit -m "feat: add enrichVcsConnectionId helper function with tests"
```

---

### Task 3: Integrate Enrichment into Template Resource

**Files:**
- Modify: `env0/resource_template.go:291-307` (create) and `env0/resource_template.go:331-345` (update)

- [ ] **Step 1: Add enrichment to resourceTemplateCreate**

In `env0/resource_template.go`, in `resourceTemplateCreate` (line 291), add the enrichment call between payload creation and the API call. The function should become:

```go
func resourceTemplateCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	request, problem := templateCreatePayloadFromParameters("", d)
	if problem != nil {
		return problem
	}

	if err := enrichVcsConnectionId(apiClient, request.GithubInstallationId, request.BitbucketClientKey, &request.VcsConnectionId); err != nil {
		return diag.FromErr(err)
	}

	template, err := apiClient.TemplateCreate(request)
	if err != nil {
		return diag.Errorf("could not create template: %v", err)
	}

	d.SetId(template.Id)

	return nil
}
```

- [ ] **Step 2: Add enrichment to resourceTemplateUpdate**

In `resourceTemplateUpdate` (line 331), add the same enrichment call:

```go
func resourceTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	request, problem := templateCreatePayloadFromParameters("", d)
	if problem != nil {
		return problem
	}

	if err := enrichVcsConnectionId(apiClient, request.GithubInstallationId, request.BitbucketClientKey, &request.VcsConnectionId); err != nil {
		return diag.FromErr(err)
	}

	_, err := apiClient.TemplateUpdate(d.Id(), request)
	if err != nil {
		return diag.Errorf("could not update template: %v", err)
	}

	return nil
}
```

- [ ] **Step 3: Verify it compiles**

Run: `go build ./...`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add env0/resource_template.go
git commit -m "feat: enrich vcs_connection_id in template create/update"
```

---

### Task 4: Integrate Enrichment into Environment Resource

**Files:**
- Modify: `env0/resource_environment.go:577-602` (create) and `env0/resource_environment.go:834-847` (update)

- [ ] **Step 1: Add enrichment to createEnvironmentWithoutTemplate**

In `env0/resource_environment.go`, in `createEnvironmentWithoutTemplate` (line 577), add the enrichment call after `templateCreatePayloadFromParameters` but before building the `EnvironmentCreateWithoutTemplate` payload:

```go
func createEnvironmentWithoutTemplate(d *schema.ResourceData, apiClient client.ApiClientInterface) (client.Environment, client.EnvironmentCreate, diag.Diagnostics) {
	templatePayload, diagError := templateCreatePayloadFromParameters("without_template_settings.0", d)
	if diagError != nil {
		return client.Environment{}, client.EnvironmentCreate{}, diagError
	}

	if err := enrichVcsConnectionId(apiClient, templatePayload.GithubInstallationId, templatePayload.BitbucketClientKey, &templatePayload.VcsConnectionId); err != nil {
		return client.Environment{}, client.EnvironmentCreate{}, diag.FromErr(err)
	}

	environmentPayload, diagError := getCreatePayload(d, apiClient, templatePayload.Type)
	if diagError != nil {
		return client.Environment{}, client.EnvironmentCreate{}, diagError
	}

	payload := client.EnvironmentCreateWithoutTemplate{
		EnvironmentCreate: environmentPayload,
		TemplateCreate:    templatePayload,
	}

	// Note: the blueprint id field of the environment is returned only during creation of a template without environment.
	// Afterward, it will be omitted from future response.
	// setEnvironmentSchema() sets the blueprint id in the resource (under "without_template_settings.0.id").
	environment, err := apiClient.EnvironmentCreateWithoutTemplate(payload)
	if err != nil {
		return client.Environment{}, client.EnvironmentCreate{}, diag.Errorf("could not create environment: %v", err)
	}

	return environment, environmentPayload, nil
}
```

- [ ] **Step 2: Add enrichment to updateTemplate**

In `updateTemplate` (line 834), add the enrichment call:

```go
func updateTemplate(d *schema.ResourceData, apiClient client.ApiClientInterface) diag.Diagnostics {
	payload, problem := templateCreatePayloadFromParameters("without_template_settings.0", d)
	if problem != nil {
		return problem
	}

	if err := enrichVcsConnectionId(apiClient, payload.GithubInstallationId, payload.BitbucketClientKey, &payload.VcsConnectionId); err != nil {
		return diag.FromErr(err)
	}

	templateId := d.Get("without_template_settings.0.id").(string)

	if _, err := apiClient.TemplateUpdate(templateId, payload); err != nil {
		return diag.Errorf("could not update template: %v", err)
	}

	return nil
}
```

- [ ] **Step 3: Verify it compiles**

Run: `go build ./...`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add env0/resource_environment.go
git commit -m "feat: enrich vcs_connection_id in environment create/update"
```

---

### Task 5: Integrate Enrichment into Custom Flow Resource

**Files:**
- Modify: `env0/resource_custom_flow.go:27-47` (create) and `env0/resource_custom_flow.go:66-83` (update)

- [ ] **Step 1: Add enrichment to resourceCustomFlowCreate**

In `env0/resource_custom_flow.go`, in `resourceCustomFlowCreate` (line 27), add between `Invalidate()` and `CustomFlowCreate`:

```go
func resourceCustomFlowCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.CustomFlowCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if err := payload.Invalidate(); err != nil {
		return diag.Errorf("invalid custom flow payload: %v", err)
	}

	if err := enrichVcsConnectionId(apiClient, payload.GithubInstallationId, payload.BitbucketClientKey, &payload.VcsConnectionId); err != nil {
		return diag.FromErr(err)
	}

	customFlow, err := apiClient.CustomFlowCreate(payload)
	if err != nil {
		return diag.Errorf("could not create custom flow: %v", err)
	}

	d.SetId(customFlow.Id)

	return nil
}
```

- [ ] **Step 2: Add enrichment to resourceCustomFlowUpdate**

In `resourceCustomFlowUpdate` (line 66):

```go
func resourceCustomFlowUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.CustomFlowCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if err := payload.Invalidate(); err != nil {
		return diag.Errorf("invalid custom flow payload: %v", err)
	}

	if err := enrichVcsConnectionId(apiClient, payload.GithubInstallationId, payload.BitbucketClientKey, &payload.VcsConnectionId); err != nil {
		return diag.FromErr(err)
	}

	if _, err := apiClient.CustomFlowUpdate(d.Id(), payload); err != nil {
		return diag.Errorf("could not update custom flow: %v", err)
	}

	return nil
}
```

- [ ] **Step 3: Verify it compiles**

Run: `go build ./...`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add env0/resource_custom_flow.go
git commit -m "feat: enrich vcs_connection_id in custom flow create/update"
```

---

### Task 6: Integrate Enrichment into Approval Policy Resource

**Files:**
- Modify: `env0/resource_approval_policy.go:27-47` (create) and `env0/resource_approval_policy.go:73-86` (update)

- [ ] **Step 1: Add enrichment to resourceApprovalPolicyCreate**

In `env0/resource_approval_policy.go`, in `resourceApprovalPolicyCreate` (line 27):

```go
func resourceApprovalPolicyCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.ApprovalPolicyCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if err := payload.Invalidate(); err != nil {
		return diag.FromErr(err)
	}

	if err := enrichVcsConnectionId(apiClient, payload.GithubInstallationId, payload.BitbucketClientKey, &payload.VcsConnectionId); err != nil {
		return diag.FromErr(err)
	}

	approvalPolicy, err := apiClient.ApprovalPolicyCreate(&payload)
	if err != nil {
		return diag.Errorf("failed to create approval policy: %v", err)
	}

	d.SetId(approvalPolicy.Id)

	return nil
}
```

- [ ] **Step 2: Add enrichment to resourceApprovalPolicyUpdate**

In `resourceApprovalPolicyUpdate` (line 73):

```go
func resourceApprovalPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.ApprovalPolicyUpdatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if err := enrichVcsConnectionId(apiClient, payload.GithubInstallationId, payload.BitbucketClientKey, &payload.VcsConnectionId); err != nil {
		return diag.FromErr(err)
	}

	if _, err := apiClient.ApprovalPolicyUpdate(&payload); err != nil {
		return diag.Errorf("failed to update approval policy: %v", err)
	}

	return nil
}
```

- [ ] **Step 3: Verify it compiles**

Run: `go build ./...`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add env0/resource_approval_policy.go
git commit -m "feat: enrich vcs_connection_id in approval policy create/update"
```

---

### Task 7: Integrate Enrichment into Module Resource

**Files:**
- Modify: `env0/resource_module.go:175-199` (create) and `env0/resource_module.go:228-249` (update)

Note: Module uses `*int` for `GithubInstallationId` and `*string` for `BitbucketClientKey` in the read model, but `*int` and `string` in the create payload and `*int` and `string` in the update payload. We need to dereference the pointer for `GithubInstallationId`.

- [ ] **Step 1: Add enrichment to resourceModuleCreate**

In `env0/resource_module.go`, in `resourceModuleCreate` (line 175), add between `Invalidate()` / module test validation and `ModuleCreate`:

```go
func resourceModuleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.ModuleCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if err := payload.Invalidate(); err != nil {
		return diag.Errorf("invalid module payload: %v", err)
	}

	if !payload.ModuleTestEnabled && (payload.RunTestsOnPullRequest || payload.OpentofuVersion != "") {
		return diag.Errorf("'run_tests_on_pull_request' and/or 'opentofu_version' may only be set if 'module_test_enabled' is enabled (set to 'true')")
	}

	ghId := 0
	if payload.GithubInstallationId != nil {
		ghId = *payload.GithubInstallationId
	}

	if err := enrichVcsConnectionId(apiClient, ghId, payload.BitbucketClientKey, &payload.VcsConnectionId); err != nil {
		return diag.FromErr(err)
	}

	module, err := apiClient.ModuleCreate(payload)
	if err != nil {
		return diag.Errorf("could not create module: %v", err)
	}

	d.SetId(module.Id)

	return nil
}
```

- [ ] **Step 2: Add enrichment to resourceModuleUpdate**

In `resourceModuleUpdate` (line 228):

```go
func resourceModuleUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.ModuleUpdatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if err := payload.Invalidate(); err != nil {
		return diag.Errorf("invalid module payload: %v", err)
	}

	if !payload.ModuleTestEnabled && (payload.RunTestsOnPullRequest || payload.OpentofuVersion != "") {
		return diag.Errorf("'run_tests_on_pull_request' and/or 'opentofu_version' may only be set if 'module_test_enabled' is enabled (set to 'true')")
	}

	ghId := 0
	if payload.GithubInstallationId != nil {
		ghId = *payload.GithubInstallationId
	}

	if err := enrichVcsConnectionId(apiClient, ghId, payload.BitbucketClientKey, &payload.VcsConnectionId); err != nil {
		return diag.FromErr(err)
	}

	if _, err := apiClient.ModuleUpdate(d.Id(), payload); err != nil {
		return diag.Errorf("could not update module: %v", err)
	}

	return nil
}
```

- [ ] **Step 3: Verify it compiles**

Run: `go build ./...`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add env0/resource_module.go
git commit -m "feat: enrich vcs_connection_id in module create/update"
```

---

### Task 8: Integrate Enrichment into Environment Discovery Resource

**Files:**
- Modify: `env0/resource_environment_discovery_configuration.go:274-346`

- [ ] **Step 1: Add enrichment to resourceEnvironmentDiscoveryConfigurationPut**

In `env0/resource_environment_discovery_configuration.go`, in `resourceEnvironmentDiscoveryConfigurationPut` (line 274), add the enrichment call between the `Invalidate()` block and the `PutEnvironmentDiscovery` API call. Insert after line 316 (after `Invalidate()`) and before line 318 (the `DiscoveryFileConfiguration == nil` check for ssh keys):

```go
	if err := putPayload.Invalidate(); err != nil {
		return diag.Errorf("invalid environment discovery payload: %v", err)
	}

	if err := enrichVcsConnectionId(apiClient, putPayload.GithubInstallationId, putPayload.BitbucketClientKey, &putPayload.VcsConnectionId); err != nil {
		return diag.FromErr(err)
	}

	if putPayload.DiscoveryFileConfiguration == nil {
```

- [ ] **Step 2: Verify it compiles**

Run: `go build ./...`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add env0/resource_environment_discovery_configuration.go
git commit -m "feat: enrich vcs_connection_id in environment discovery put"
```

---

### Task 9: Update Existing Resource Tests to Mock VcsConnections

**Files:**
- Modify: `env0/resource_custom_flow_test.go`
- Modify: `env0/resource_template_test.go`
- Modify: `env0/resource_approval_policy_test.go`
- Modify: `env0/resource_module_test.go`
- Modify: `env0/resource_environment_test.go`
- Modify: `env0/resource_environment_discovery_configuration_test.go`

The enrichment function calls `VcsConnections()` whenever `github_installation_id != 0` or `bitbucket_client_key != ""`. Existing tests that set these fields now need a mock expectation for `VcsConnections()`.

The simplest approach: return an empty list from `VcsConnections()`. Since no match is found and `enrichVcsConnectionId` returns an error, the test would fail. Instead, return a list containing a connection that matches the test's `github_installation_id` or `bitbucket_client_key`.

For each test file, find test cases that use `github_installation_id` (non-zero) or `bitbucket_client_key` (non-empty) in their config, and add a `VcsConnections()` mock expectation that returns a matching connection.

**Strategy:** Add a shared mock expectation pattern. For each test that sets `github_installation_id` to a non-zero value, add:

```go
mock.EXPECT().VcsConnections().AnyTimes().Return([]client.VcsConnection{
	{Id: "vcs-conn-enriched", GithubInstallationId: <the_github_installation_id_value>},
}, nil)
```

And update the expected create/update payload to include `VcsConnectionId: "vcs-conn-enriched"`.

For tests that only use `vcs_connection_id` (already set), or use neither `github_installation_id` nor `bitbucket_client_key`, no changes are needed since the enrichment is a no-op.

- [ ] **Step 1: Run existing tests to identify failures**

Run: `go test ./env0/ -v -count=1 2>&1 | head -200`

Identify which tests fail due to unexpected `VcsConnections()` calls.

- [ ] **Step 2: Fix each failing test file**

For each failing test, add the `VcsConnections()` mock expectation and update the expected payload's `VcsConnectionId` field. The exact changes depend on which tests fail — fix each one based on the test's `github_installation_id` or `bitbucket_client_key` value.

Pattern for fixing: wherever `mock.EXPECT().CustomFlowCreate(createPayload)` (or similar) is called and the test uses `github_installation_id`, add:
1. `mock.EXPECT().VcsConnections().AnyTimes().Return(...)`
2. Update `createPayload.VcsConnectionId = "vcs-conn-enriched"` in the expected payload

- [ ] **Step 3: Run all tests to verify they pass**

Run: `go test ./env0/ -v -count=1`
Expected: All tests PASS

- [ ] **Step 4: Commit**

```bash
git add env0/*_test.go
git commit -m "test: update resource tests to mock VcsConnections for enrichment"
```

---

### Task 10: Final Verification

- [ ] **Step 1: Run full test suite**

Run: `go test ./... -count=1`
Expected: All tests PASS

- [ ] **Step 2: Run linter if available**

Run: `go vet ./...`
Expected: No errors

- [ ] **Step 3: Verify build**

Run: `go build ./...`
Expected: No errors
