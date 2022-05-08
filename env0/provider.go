package env0

import (
	"context"

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
			},
			DataSourcesMap: map[string]*schema.Resource{
				"env0_organization":            dataOrganization(),
				"env0_project":                 dataProject(),
				"env0_project_policy":          dataPolicy(),
				"env0_configuration_variable":  dataConfigurationVariable(),
				"env0_template":                dataTemplate(),
				"env0_templates":               dataTemplates(),
				"env0_ssh_key":                 dataSshKey(),
				"env0_aws_cost_credentials":    dataCostCredentials(string(client.AwsCostCredentialsType)),
				"env0_azure_cost_credentials":  dataCostCredentials(string(client.AzureCostCredentialsType)),
				"env0_google_cost_credentials": dataCostCredentials(string(client.GoogleCostCredentialsType)),
				"env0_aws_credentials":         dataAwsCredentials(),
				"env0_gcp_credentials":         dataGcpCredentials(),
				"env0_azure_credentials":       dataAzureCredentials(),
				"env0_team":                    dataTeam(),
				"env0_teams":                   dataTeams(),
				"env0_environment":             dataEnvironment(),
				"env0_workflow_triggers":       dataWorkflowTriggers(),
				"env0_notification":            dataNotification(),
				"env0_module":                  dataModule(),
				"env0_git_token":               dataGitToken(),
				"env0_api_key":                 dataApiKey(),
				"env0_agents":                  dataAgents(),
				"env0_user":                    dataUser(),
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
				"env0_environment_scheduling":               resourceEnvironmentScheduling(),
				"env0_environment_drift_detection":          resourceDriftDetection(),
				"env0_notification":                         resourceNotification(),
				"env0_notification_project_assignment":      resourceNotificationProjectAssignment(),
				"env0_module":                               resourceModule(),
				"env0_git_token":                            resourceGitToken(),
				"env0_api_key":                              resourceApiKey(),
				"env0_organization_policy":                  resourceOrganizationPolicy(),
				"env0_agent_project_assignment":             resourceAgentProjectAssignment(),
			},
		}

		provider.ConfigureContextFunc = configureProvider(version, provider)
		return provider
	}
}

func configureProvider(version string, p *schema.Provider) schema.ConfigureContextFunc {
	userAgent := p.UserAgent("terraform-provider-env0", version)

	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		httpClient, err := http.NewHttpClient(http.HttpClientConfig{
			ApiKey:      d.Get("api_key").(string),
			ApiSecret:   d.Get("api_secret").(string),
			ApiEndpoint: d.Get("api_endpoint").(string),
			UserAgent:   userAgent,
			RestClient:  resty.New(),
		})
		if err != nil {
			return nil, diag.Diagnostics{diag.Diagnostic{Severity: diag.Error, Summary: err.Error()}}
		}

		return client.NewApiClient(httpClient), nil
	}
}
