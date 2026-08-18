package env0

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/utils"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type testRestyClientSuite struct {
	suite.Suite
	client *resty.Client
	url    string
}

// TestMain sets the shared env vars once, before any test runs. Doing this per-call in runUnitTest
// would be a data race now that top-level tests run with t.Parallel().
func TestMain(m *testing.M) {
	os.Setenv("TF_ACC", "1")
	os.Setenv("ENV0_API_KEY", "value")
	os.Setenv("ENV0_API_SECRET", "value")
	os.Exit(m.Run())
}

func runUnitTest(t *testing.T, testCase resource.TestCase, mockFunc func(mockFunc *client.MockApiClientInterface)) {
	t.Helper()

	testPattern := os.Getenv("TEST_PATTERN")
	if testPattern != "" && !strings.Contains(t.Name(), testPattern) {
		t.SkipNow()

		return
	}

	testReporter := utils.TestReporter{T: t}
	ctrl := gomock.NewController(&testReporter)

	apiClientMock := client.NewMockApiClientInterface(ctrl)
	mockFunc(apiClientMock)

	testCase.ProviderFactories = map[string]func() (*schema.Provider, error){
		//nolint:all // tests
		"env0": func() (*schema.Provider, error) {
			provider := Provider("")()
			provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
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

	diags := Provider("TEST")().Configure(context.Background(), &terraform.ResourceConfig{})
	testExpectedProviderError(t, diags, expectedKey)
}

func testMissingConfig(t *testing.T, config map[string]any, expectedKey string) {
	diags := Provider("TEST")().Configure(context.Background(), terraform.NewResourceConfigRaw(config))
	testExpectedProviderError(t, diags, expectedKey)
}

func TestMissingConfigurations(t *testing.T) {
	expectedApiKeyConfig := "api_key"
	expectedApiSecretConfig := "api_secret"

	configTestCases := map[string]map[string]any{
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

func (suite *testRestyClientSuite) SetupTest() {
	httpmock.Reset()
}

func (suite *testRestyClientSuite) SetupSuite() {
	httpmock.ActivateNonDefault(suite.client.GetClient())
}

func (suite *testRestyClientSuite) TearDownAllSuite() {
	httpmock.Deactivate()
}

func (suite *testRestyClientSuite) TestOkResponse() {
	t := suite.T()

	httpmock.RegisterResponder("GET", suite.url, httpmock.NewStringResponder(http.StatusOK, "OK"))

	res, err := suite.client.R().Get(suite.url)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.StatusCode())
		assert.Equal(t, "OK", res.String())
	}

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func (suite *testRestyClientSuite) Test4xxResponse() {
	t := suite.T()

	httpmock.RegisterResponder("GET", suite.url, httpmock.NewStringResponder(http.StatusBadRequest, "BAD"))

	res, err := suite.client.R().Get(suite.url)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
		assert.Equal(t, "BAD", res.String())
	}

	// Should be called once - no retries.
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

// REMOVED - takes too long to run.
// func (suite *testRestyClientSuite) Test5xxResponse() {
// 	t := suite.T()

// 	httpmock.RegisterResponder("GET", suite.url, httpmock.NewStringResponder(http.StatusInternalServerError, "BAD"))

// 	res, err := suite.client.R().Get(suite.url)

// 	if assert.NoError(t, err) {
// 		assert.Equal(t, http.StatusInternalServerError, res.StatusCode())
// 		assert.Equal(t, "BAD", res.String())
// 	}

// 	// Should be called multiple times - retries.
// 	assert.Equal(t, 11, httpmock.GetTotalCallCount())
// }

func TestRestyClientSuite(t *testing.T) {
	s := &testRestyClientSuite{
		client: createRestyClient(context.Background()),
		url:    "http://fake.env0.com/fake",
	}
	suite.Run(t, s)
}

// createRetryTestClient returns a client with the production retry conditions but a negligible
// backoff, so asserting on attempt counts doesn't pay the real ladder's wall time.
func createRetryTestClient(t *testing.T) *resty.Client {
	t.Helper()

	return createRestyClient(context.Background()).
		SetRetryWaitTime(time.Millisecond).
		SetRetryMaxWaitTime(time.Millisecond)
}

// An empty list is a legitimate answer for most list endpoints. The integration-test-only
// retry that covers read-after-write lag must stop after emptyListMaxAttempts, or every
// genuinely-empty list pays the full retry ladder.
func TestRestyClientEmptyListRetryIsCapped(t *testing.T) {
	t.Setenv("INTEGRATION_TESTS", "1")

	client := createRetryTestClient(t)
	url := "http://fake.env0.com/empty-list"

	httpmock.ActivateNonDefault(client.GetClient())

	defer httpmock.Deactivate()

	httpmock.Reset()
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(http.StatusOK, "[]"))

	res, err := client.R().Get(url)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.StatusCode())
		assert.Equal(t, "[]", res.String())
	}

	assert.Equal(t, emptyListMaxAttempts, httpmock.GetTotalCallCount())
}

// Some endpoints answer 404 by design, so the integration-test retry that covers database
// eventual consistency must stop after notFoundMaxAttempts.
func TestRestyClientNotFoundRetryIsCapped(t *testing.T) {
	t.Setenv("INTEGRATION_TESTS", "1")

	client := createRetryTestClient(t)
	url := "http://fake.env0.com/not-found"

	httpmock.ActivateNonDefault(client.GetClient())

	defer httpmock.Deactivate()

	httpmock.Reset()
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(http.StatusNotFound, "NOT FOUND"))

	res, err := client.R().Get(url)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNotFound, res.StatusCode())
	}

	assert.Equal(t, notFoundMaxAttempts, httpmock.GetTotalCallCount())
}

// Outside the integration tests a 404 is never retried.
func TestRestyClientNotFoundIsNotRetried(t *testing.T) {
	t.Setenv("INTEGRATION_TESTS", "")

	client := createRetryTestClient(t)
	url := "http://fake.env0.com/not-found-no-integration"

	httpmock.ActivateNonDefault(client.GetClient())

	defer httpmock.Deactivate()

	httpmock.Reset()
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(http.StatusNotFound, "NOT FOUND"))

	res, err := client.R().Get(url)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNotFound, res.StatusCode())
	}

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

// A 429 is the failure this provider has to survive: the whole retry ladder is spent on it.
func TestRestyClientRateLimitIsRetried(t *testing.T) {
	client := createRetryTestClient(t)
	url := "http://fake.env0.com/rate-limited"

	httpmock.ActivateNonDefault(client.GetClient())

	defer httpmock.Deactivate()

	httpmock.Reset()
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(http.StatusTooManyRequests, "TOO MANY REQUESTS"))

	res, err := client.R().Get(url)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusTooManyRequests, res.StatusCode())
	}

	assert.Equal(t, retryCount+1, httpmock.GetTotalCallCount())
}

func newRetryTestResponse(statusCode int, attempt int, retryAfter string) *resty.Response {
	header := http.Header{}
	if retryAfter != "" {
		header.Set("Retry-After", retryAfter)
	}

	return &resty.Response{
		Request:     &resty.Request{Attempt: attempt},
		RawResponse: &http.Response{StatusCode: statusCode, Header: header},
	}
}

// A 429 comes from a minute-sized window, so its ladder has to start higher and climb further
// than the one used for a 5xx, and it has to follow Retry-After when the response supplies it.
func TestRetryWaitTimeFor(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		attempt    int
		retryAfter string
		min        time.Duration
		max        time.Duration
	}{
		{"5xx first retry", http.StatusInternalServerError, 1, "", retryWaitTime / 2, retryWaitTime},
		{"5xx third retry", http.StatusInternalServerError, 3, "", retryWaitTime * 2, retryWaitTime * 4},
		{"5xx is capped", http.StatusInternalServerError, retryCount, "", retryMaxWaitTime / 2, retryMaxWaitTime},
		{"429 first retry", http.StatusTooManyRequests, 1, "", rateLimitWaitTime / 2, rateLimitWaitTime},
		{"429 is capped", http.StatusTooManyRequests, retryCount, "", rateLimitMaxWaitTime / 2, rateLimitMaxWaitTime},
		{"429 with retry after seconds", http.StatusTooManyRequests, 1, "8", time.Second * 8, time.Second * 10},
		{"429 with unparsable retry after", http.StatusTooManyRequests, 1, "soon", rateLimitWaitTime / 2, rateLimitWaitTime},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			waitTime := retryWaitTimeFor(newRetryTestResponse(tt.statusCode, tt.attempt, tt.retryAfter))

			assert.GreaterOrEqual(t, waitTime, tt.min)
			assert.LessOrEqual(t, waitTime, tt.max)
		})
	}
}

// The 429 ladder must not collapse to the 5xx one: a 429 always waits longer than a 5xx would
// have waited on the same attempt.
func TestRetryWaitTimeForRateLimitWaitsLongerThan5xx(t *testing.T) {
	for attempt := 1; attempt <= retryCount; attempt++ {
		rateLimited := retryWaitTimeFor(newRetryTestResponse(http.StatusTooManyRequests, attempt, ""))
		serverError := retryWaitTimeFor(newRetryTestResponse(http.StatusInternalServerError, attempt, ""))

		assert.Greater(t, rateLimited, serverError, "attempt %d", attempt)
	}
}

// Requests blocked in the same window all get the same Retry-After. Retrying them all the moment
// it expires is what got the provider blocked by the WAF, so the wait is jittered.
func TestRetryWaitTimeForRateLimitIsJittered(t *testing.T) {
	waitTimes := map[time.Duration]bool{}

	for range 20 {
		waitTimes[retryWaitTimeFor(newRetryTestResponse(http.StatusTooManyRequests, 1, "10"))] = true
	}

	assert.Greater(t, len(waitTimes), 1)
}
