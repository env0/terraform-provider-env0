package env0

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func Provider(version string) plugin.ProviderFunc {
	return func() *schema.Provider {
		apiKeyEnv := "ENV0_API_KEY"
		apiSecretEnv := "ENV0_API_SECRET"
		apiOrganizationIdEnv := "ENV0_ORGANIZATION_ID"

		// version "TEST" is used for acceptance testing the provider with env variables.
		// must use different env variables to prevent race conditions with other acceptance tests that run in parallel.
		if version == "TEST" {
			version = ""
			apiKeyEnv = "ENV0_API_KEY_TEST"
			apiSecretEnv = "ENV0_API_SECRET_TEST"
		}

		provider := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"api_endpoint": {
					Type:        schema.TypeString,
					Description: "env0 API endpoint. This can also be set via the ENV0_API_ENDPOINT environment variable, and is usually used for testing purposes.",
					DefaultFunc: schema.EnvDefaultFunc("ENV0_API_ENDPOINT", "https://api.env0.com/"),
					Optional:    true,
				},
				"api_key": {
					Type:        schema.TypeString,
					Description: "env0 API key. This can also be set via the ENV0_API_KEY environment variable.",
					DefaultFunc: schema.EnvDefaultFunc(apiKeyEnv, nil),
					Required:    true,
					Sensitive:   true,
				},
				"api_secret": {
					Type:        schema.TypeString,
					Description: "env0 API secret. This can also be set via the ENV0_API_SECRET environment variable.",
					DefaultFunc: schema.EnvDefaultFunc(apiSecretEnv, nil),
					Required:    true,
					Sensitive:   true,
				},
				"organization_id": {
					Type:        schema.TypeString,
					Description: "when the API key is associated with multiple organizations, this field is required. If an API key has one organization, this field is ignored. This can also be set via the ENV0_ORGANIZATION_ID environment variable.",
					DefaultFunc: schema.EnvDefaultFunc(apiOrganizationIdEnv, nil),
					Optional:    true,
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"env0_organization":              dataOrganization(),
				"env0_project":                   dataProject(),
				"env0_project_policy":            dataPolicy(),
				"env0_configuration_variable":    dataConfigurationVariable(),
				"env0_template":                  dataTemplate(),
				"env0_templates":                 dataTemplates(),
				"env0_ssh_key":                   dataSshKey(),
				"env0_aws_cost_credentials":      dataCredentials(AWS_COST_TYPE),
				"env0_azure_cost_credentials":    dataCredentials(AZURE_COST_TYPE),
				"env0_google_cost_credentials":   dataCredentials(GCP_COST_TYPE),
				"env0_aws_credentials":           dataCredentials(AWS_TYPE),
				"env0_gcp_credentials":           dataCredentials(GCP_TYPE),
				"env0_azure_credentials":         dataCredentials(AZURE_TYPE),
				"env0_team":                      dataTeam(),
				"env0_teams":                     dataTeams(),
				"env0_environment":               dataEnvironment(),
				"env0_workflow_triggers":         dataWorkflowTriggers(),
				"env0_notification":              dataNotification(),
				"env0_notifications":             dataNotifications(),
				"env0_module":                    dataModule(),
				"env0_git_token":                 dataGitToken(),
				"env0_api_key":                   dataApiKey(),
				"env0_agents":                    dataAgents(),
				"env0_user":                      dataUser(),
				"env0_cloud_credentials":         dataCloudCredentials(),
				"env0_project_cloud_credentials": dataProjectCloudCredentials(),
				"env0_source_code_variables":     dataSourceCodeVariables(),
				"env0_ip_ranges":                 dataIpRanges(),
				"env0_agent_values":              dataAgentValues(),
				"env0_custom_role":               dataCustomRole(),
				"env0_custom_roles":              dataCustomRoles(),
				"env0_gpg_key":                   dataGpgKey(),
				"env0_provider":                  dataProvider(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"env0_project":                              resourceProject(),
				"env0_project_policy":                       resourcePolicy(),
				"env0_configuration_variable":               resourceConfigurationVariable(),
				"env0_template":                             resourceTemplate(),
				"env0_ssh_key":                              resourceSshKey(),
				"env0_aws_credentials":                      resourceAwsCredentials(),
				"env0_aws_cost_credentials":                 resourceCostCredentials("aws"),
				"env0_azure_cost_credentials":               resourceCostCredentials("azure"),
				"env0_gcp_cost_credentials":                 resourceCostCredentials("google"),
				"env0_gcp_credentials":                      resourceGcpCredentials(),
				"env0_azure_credentials":                    resourceAzureCredentials(),
				"env0_template_project_assignment":          resourceTemplateProjectAssignment(),
				"env0_cloud_credentials_project_assignment": resourceCloudCredentialsProjectAssignment(),
				"env0_cost_credentials_project_assignment":  resourceCostCredentialsProjectAssignment(),
				"env0_team_project_assignment":              resourceTeamProjectAssignment(),
				"env0_team":                                 resourceTeam(),
				"env0_environment":                          resourceEnvironment(),
				"env0_workflow_triggers":                    resourceWorkflowTriggers(),
				"env0_workflow_trigger":                     resourceWorkflowTrigger(),
				"env0_environment_scheduling":               resourceEnvironmentScheduling(),
				"env0_environment_drift_detection":          resourceDriftDetection(),
				"env0_notification":                         resourceNotification(),
				"env0_notification_project_assignment":      resourceNotificationProjectAssignment(),
				"env0_module":                               resourceModule(),
				"env0_git_token":                            resourceGitToken(),
				"env0_api_key":                              resourceApiKey(),
				"env0_organization_policy":                  resourceOrganizationPolicy(),
				"env0_agent_project_assignment":             resourceAgentProjectAssignment(),
				"env0_user_team_assignment":                 resourceUserTeamAssignment(),
				"env0_user_project_assignment":              resourceUserProjectAssignment(),
				"env0_custom_role":                          resourceCustomRole(),
				"env0_user_organization_assignment":         resourceUserOrganizationAssignment(),
				"env0_custom_flow":                          resourceCustomFlow(),
				"env0_custom_flow_assignment":               resourceCustomFlowAssignment(),
				"env0_environment_state_access":             resourceEnvironmentStateAccess(),
				"env0_gpg_key":                              resourceGpgKey(),
				"env0_provider":                             resourceProvider(),
				"env0_user_environment_assignment":          resourceUserEnvironmentAssignment(),
				"env0_team_environment_assignment":          resourceTeamEnvironmentAssignment(),
			},
		}

		provider.ConfigureContextFunc = configureProvider(version, provider)
		return provider
	}
}

func configureProvider(version string, p *schema.Provider) schema.ConfigureContextFunc {
	userAgent := p.UserAgent("terraform-provider-env0", version)

	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		restyClient := resty.New()

		isIntegrationTest := false

		if os.Getenv("INTEGRATION_TESTS") == "1" {
			isIntegrationTest = true
		}

		restyClient.
			SetRetryCount(5).
			SetRetryWaitTime(time.Second).
			SetRetryMaxWaitTime(time.Second * 5).
			AddRetryCondition(func(r *resty.Response, err error) bool {
				if r == nil {
					// No response. Possiblly a networking issue (E.g. DNS lookup failure).
					log.Printf("[WARN] No response, retrying request: %s %s", r.Request.Method, r.Request.URL)
					return true
				}

				// When running integration tests 404 may occur due to "database eventual consistency".
				// Retry when there's a 5xx error. Otherwise do not retry.
				if r.StatusCode() >= 500 || isIntegrationTest && r.StatusCode() == 404 {
					log.Printf("[WARN] Received %d status code, retrying request: %s %s", r.StatusCode(), r.Request.Method, r.Request.URL)
					return true
				}

				if r.StatusCode() == 200 && isIntegrationTest && r.String() == "[]" {
					log.Printf("[WARN] Received an empty list for an integration test, retrying request: %s %s", r.Request.Method, r.Request.URL)
					return true
				}

				return false
			})

		httpClient, err := http.NewHttpClient(http.HttpClientConfig{
			ApiKey:      d.Get("api_key").(string),
			ApiSecret:   d.Get("api_secret").(string),
			ApiEndpoint: d.Get("api_endpoint").(string),
			UserAgent:   userAgent,
			RestClient:  restyClient,
		})
		if err != nil {
			return nil, diag.Diagnostics{diag.Diagnostic{Severity: diag.Error, Summary: err.Error()}}
		}

		apiClient := client.NewApiClient(httpClient, d.Get("organization_id").(string))

		// organizations fetched to cache Auth0 API response.
		if _, err := apiClient.OrganizationId(); err != nil {
			return nil, diag.Diagnostics{diag.Diagnostic{Severity: diag.Error, Summary: err.Error()}}
		}
		return apiClient, nil
	}
}
