package env0

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/utils"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func runUnitTest(t *testing.T, testCase resource.TestCase, mockFunc func(mockFunc *client.MockApiClientInterface)) {
	os.Setenv("TF_ACC", "1")
	os.Setenv("ENV0_API_KEY", "value")
	os.Setenv("ENV0_API_SECRET", "value")

	testPattern := os.Getenv("TEST_PATTERN")
	if testPattern != "" && !strings.Contains(t.Name(), testPattern) {
		t.SkipNow()
		return
	}

	testReporter := utils.TestReporter{T: t}
	ctrl := gomock.NewController(&testReporter)
	defer ctrl.Finish()

	apiClientMock := client.NewMockApiClientInterface(ctrl)
	mockFunc(apiClientMock)

	testCase.ProviderFactories = map[string]func() (*schema.Provider, error){
		"env0": func() (*schema.Provider, error) {
			provider := Provider("")()
			provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
				return apiClientMock, nil
			}
			return provider, nil
		},
	}
	testCase.PreventPostDestroyRefresh = true
	resource.ParallelTest(&testReporter, testCase)
}

func TestProvider(t *testing.T) {
	if err := Provider("")().InternalValidate(); err != nil {
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
		t.Fatalf("Error wasn't received, expected: %s", expectedError)
	}
}

func testMissingEnvVar(t *testing.T, envVars map[string]string, expectedKey string) {
	for key, value := range envVars {
		os.Setenv(key, value)
		defer os.Setenv(key, "")
	}

	diags := Provider("TEST")().Validate(&terraform.ResourceConfig{})
	testExpectedProviderError(t, diags, expectedKey)
}

func testMissingConfig(t *testing.T, config map[string]interface{}, expectedKey string) {
	diags := Provider("TEST")().Validate(terraform.NewResourceConfigRaw(config))
	testExpectedProviderError(t, diags, expectedKey)
}

func TestMissingConfigurations(t *testing.T) {
	expectedApiKeyConfig := "api_key"
	expectedApiSecretConfig := "api_secret"

	configTestCases := map[string]map[string]interface{}{
		expectedApiKeyConfig: {
			"api_secret": "value",
		},
		expectedApiSecretConfig: {
			"api_key": "value",
		},
	}

	for expectedError, config := range configTestCases {
		testMissingConfig(t, config, expectedError)
	}

	envVarsTestCases := map[string]map[string]string{
		expectedApiKeyConfig: {
			"ENV0_API_SECRET_TEST": "value",
		},
		expectedApiSecretConfig: {
			"ENV0_API_KEY_TEST": "value",
		},
	}

	for expectedError, envVars := range envVarsTestCases {
		testMissingEnvVar(t, envVars, expectedError)
	}
}
