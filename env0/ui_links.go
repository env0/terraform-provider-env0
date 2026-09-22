package env0

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
)

// uiBaseUrl maps an env0 API endpoint to the matching web UI base URL. The prod API serves the UI
// from app.env0.com; any other env0 endpoint of the form api-<host> serves it from <host> on the
// same scheme. An unrecognized endpoint (e.g. self hosted) returns an empty string.
// Note: tests/harness.go has its own environmentUiUrl with a looser api->app rule; it cannot import
// this one (package main). Keep the two in mind when changing the mapping.
func uiBaseUrl(apiEndpoint string) string {
	parsed, err := url.Parse(apiEndpoint)
	if err != nil || parsed.Host == "" || parsed.Scheme == "" {
		return ""
	}

	if parsed.Host == "api.env0.com" {
		return "https://app.env0.com"
	}

	host, found := strings.CutPrefix(parsed.Host, "api-")
	if !found || strings.HasPrefix(host, ".") || !strings.HasSuffix(host, ".env0.com") {
		return ""
	}

	return fmt.Sprintf("%s://%s", parsed.Scheme, host)
}

// deploymentUrl returns the web UI link for a deployment, the page where a deployment waiting for
// approval is approved or cancelled. It returns an empty string when the UI host is unknown.
func deploymentUrl(apiClient client.ApiClientInterface, projectId string, environmentId string, deploymentId string) string {
	baseUrl := uiBaseUrl(apiClient.ApiEndpoint())
	if baseUrl == "" {
		return ""
	}

	deploymentUrl := fmt.Sprintf("%s/p/%s/environments/%s/deployments/%s", baseUrl, projectId, environmentId, deploymentId)

	if organizationId, err := apiClient.OrganizationId(); err == nil && organizationId != "" {
		deploymentUrl += "?organizationId=" + organizationId
	}

	return deploymentUrl
}
