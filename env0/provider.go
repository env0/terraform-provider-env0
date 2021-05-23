package env0

import (
	"errors"
	"github.com/env0/terraform-provider-env0/client/http"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_endpoint": {
				Type:        schema.TypeString,
				Description: "override api endpoint (used for testing)",
				DefaultFunc: schema.EnvDefaultFunc("ENV0_API_ENDPOINT", "https://api.env0.com/"),
				Optional:    true,
				Sensitive:   false,
			},
			"api_key": {
				Type:        schema.TypeString,
				Description: "env0 api key (https://docs.env0.com/reference#authentication)",
				DefaultFunc: schema.EnvDefaultFunc("ENV0_API_KEY", nil),
				Optional:    true,
				Sensitive:   true,
			},
			"api_secret": {
				Type:        schema.TypeString,
				Description: "env0 api key secret",
				DefaultFunc: schema.EnvDefaultFunc("ENV0_API_SECRET", nil),
				Optional:    true,
				Sensitive:   true,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"env0_organization":           dataOrganization(),
			"env0_project":                dataProject(),
			"env0_configuration_variable": dataConfigurationVariable(),
			"env0_template":               dataTemplate(),
			"env0_ssh_key":                dataSshKey(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"env0_project":                resourceProject(),
			"env0_configuration_variable": resourceConfigurationVariable(),
			"env0_template":               resourceTemplate(),
			"env0_ssh_key":                resourceSshKey(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	apiKey, ok := d.GetOk("api_key")
	if !ok {
		return nil, errors.New("either api_key must be provided, or ENV0_API_KEY environment variable set")
	}
	apiSecret, ok := d.GetOk("api_secret")
	if !ok {
		return nil, errors.New("either api_secret must be provided or ENV0_API_SECRET environment variable set")
	}

	httpClient, err := http.NewHttpClient(apiKey.(string), apiSecret.(string))
	if err != nil {
		return nil, err
	}

	return client.NewApiClient(httpClient), nil
}
