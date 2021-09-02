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
		provider := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"api_endpoint": {
					Type:        schema.TypeString,
					Description: "override api endpoint (used for testing)",
					DefaultFunc: schema.EnvDefaultFunc("ENV0_API_ENDPOINT", "https://api.env0.com/"),
					Optional:    true,
				},
				"api_key": {
					Type:        schema.TypeString,
					Description: "env0 api key (https://docs.env0.com/reference#authentication)",
					DefaultFunc: schema.EnvDefaultFunc("ENV0_API_KEY", nil),
					Required:    true,
					Sensitive:   true,
				},
				"api_secret": {
					Type:        schema.TypeString,
					Description: "env0 api key secret",
					DefaultFunc: schema.EnvDefaultFunc("ENV0_API_SECRET", nil),
					Required:    true,
					Sensitive:   true,
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"env0_organization":           dataOrganization(),
				"env0_project":                dataProject(),
				"env0_configuration_variable": dataConfigurationVariable(),
				"env0_template":               dataTemplate(),
				"env0_ssh_key":                dataSshKey(),
				"env0_aws_credentials":        dataAwsCredentials(),
				"env0_team":                   dataTeam(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"env0_project":                              resourceProject(),
				"env0_configuration_variable":               resourceConfigurationVariable(),
				"env0_template":                             resourceTemplate(),
				"env0_ssh_key":                              resourceSshKey(),
				"env0_aws_credentials":                      resourceAwsCredentials(),
				"env0_template_project_assignment":          resourceTemplateProjectAssignment(),
				"env0_cloud_credentials_project_assignment": resourceCloudCredentialsProjectAssignment(),
				"env0_team_project_assignment":              resourceTeamProjectAssignment(),
				"env0_team":                                 resourceTeam(),
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
