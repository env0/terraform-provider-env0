package env0

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
)

// uiBaseUrl maps an env0 API endpoint to the matching web UI base URL. The mapping mirrors the
// FRONTEND_BASE_URL convention of the env0 stages: prod serves the UI from app.env0.com and every
// other stage from <stage>.dev.env0.com, where only the dev and bors stages are behind https.
// An unrecognized endpoint (e.g. self hosted) returns an empty string.
func uiBaseUrl(apiEndpoint string) string {
	parsed, err := url.Parse(strings.TrimSuffix(apiEndpoint, "/"))
	if err != nil || parsed.Host == "" {
		return ""
	}

	if parsed.Host == "api.env0.com" {
		return "https://app.env0.com"
	}

	stage, found := strings.CutPrefix(parsed.Host, "api-")
	stage, foundSuffix := strings.CutSuffix(stage, ".dev.env0.com")

	if !found || !foundSuffix || stage == "" {
		return ""
	}

	scheme := "http"
	if stage == "dev" || stage == "bors" {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s.dev.env0.com", scheme, stage)
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
