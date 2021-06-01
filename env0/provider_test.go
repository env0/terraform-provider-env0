package env0

import (
	"context"
	"fmt"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/utils"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"strings"
	"testing"
)

var (
	apiClientMock *client.MockApiClientInterface
	ctrl          *gomock.Controller
)

var testUnitProviders = map[string]func() (*schema.Provider, error){
	"env0": func() (*schema.Provider, error) {
		provider := Provider()
		provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			return apiClientMock, nil
		}
		return provider, nil
	},
}

func runUnitTest(t *testing.T, testCase resource.TestCase, mockFunc func(mockFunc *client.MockApiClientInterface)) {
	testReporter := utils.TestReporter{T: t}

	ctrl = gomock.NewController(&testReporter)

	apiClientMock = client.NewMockApiClientInterface(ctrl)
	mockFunc(apiClientMock)

	testCase.ProviderFactories = testUnitProviders
	resource.UnitTest(&testReporter, testCase)

	ctrl.Finish()
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testExpectedProviderError(t *testing.T, diags diag.Diagnostics, expectedKey string) {
	expectedError := fmt.Sprintf("The argument \"%s\" is required, but no definition was found.", expectedKey)
	var errorDetail string
	for _, diag := range diags {
		if strings.Contains(diag.Detail, expectedError) {
			errorDetail = diag.Detail
		}
	}

	if errorDetail == "" {
		t.Fatalf("Error wasn't recieved, expected: %s", expectedError)
	}
}

func testMissingEnvVar(t *testing.T, envVars map[string]string, expectedKey string) {
	for key, value := range envVars {
		os.Setenv(key, value)
		defer os.Setenv(key, "")
	}

	diags := Provider().Validate(&terraform.ResourceConfig{})
	testExpectedProviderError(t, diags, expectedKey)
}

func testMissingConfig(t *testing.T, config map[string]interface{}, expectedKey string) {
	diags := Provider().Validate(terraform.NewResourceConfigRaw(config))
	testExpectedProviderError(t, diags, expectedKey)
}

func TestMissingConfigurations(t *testing.T) {
	expectedApiKeyConfig := "api_key"
	expectedApiSecretConfig := "api_secret"

	configTestCases := map[string]map[string]interface{}{
		expectedApiKeyConfig: map[string]interface{}{
			"api_secret": "value",
		},
		expectedApiSecretConfig: map[string]interface{}{
			"api_key": "value",
		},
	}

	for expectedError, config := range configTestCases {
		testMissingConfig(t, config, expectedError)
	}

	envVarsTestCases := map[string]map[string]string{
		expectedApiKeyConfig: map[string]string{
			"ENV0_API_SECRET": "value",
		},
		expectedApiSecretConfig: map[string]string{
			"ENV0_API_KEY": "value",
		},
	}

	for expectedError, envVars := range envVarsTestCases {
		testMissingEnvVar(t, envVars, expectedError)
	}
}
