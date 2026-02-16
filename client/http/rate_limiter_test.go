package http_test

import (
	"time"

	httpModule "github.com/env0/terraform-provider-env0/client/http"
	"github.com/env0/terraform-provider-env0/client/http/ratelimiter"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SlidingWindow Rate Limiter", func() {
	const (
		BaseUrl         = "https://fake.env0.com"
		ApiKey          = "TEST_KEY"
		ApiSecret       = "TEST_SECRET"
		UserAgent       = "test-agent"
		TestEndpoint    = "/test-endpoint"
		SuccessResponse = "success"
	)

	var (
		restClient *resty.Client
		httpClient *httpModule.HttpClient
	)

	BeforeEach(func() {
		// Set up a new REST client for each test
		restClient = resty.New()
		httpmock.ActivateNonDefault(restClient.GetClient())

		// Register a responder for all requests to the test endpoint
		httpmock.RegisterResponder("GET", BaseUrl+TestEndpoint,
			httpmock.NewStringResponder(200, SuccessResponse))
	})

	AfterEach(func() {
		httpmock.DeactivateAndReset()
	})

	createClient := func(maxRequests int, window time.Duration) *httpModule.HttpClient {
		config := httpModule.HttpClientConfig{
			ApiKey:      ApiKey,
			ApiSecret:   ApiSecret,
			ApiEndpoint: BaseUrl,
			UserAgent:   UserAgent,
			RestClient:  restClient,
			RateLimiter: ratelimiter.NewSlidingWindowLimiter(maxRequests, window),
		}
		client, err := httpModule.NewHttpClient(config)
		Expect(err).To(BeNil())

		return client
	}

	makeRequest := func(client *httpModule.HttpClient) {
		var response string

		err := client.Get(TestEndpoint, nil, &response)
		Expect(err).To(BeNil())
		Expect(response).To(Equal(SuccessResponse))
	}

	Context("with client rate limiting tests", func() {
		// These tests verify our HTTP client's rate limiting behavior
		It("should allow multiple requests up to the limit", func() {
			const maxConcurrentRequests = 10

			httpClient = createClient(maxConcurrentRequests, 100*time.Millisecond)

			// Make a series of requests that should all succeed immediately
			for range maxConcurrentRequests {
				go makeRequest(httpClient)
			}

			// Verify all requests were made successfully
			time.Sleep(5 * time.Millisecond)

			callCount := httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(maxConcurrentRequests))
		})

		It("should handle concurrent requests with rate limiting", func() {
			const maxConcurrentRequests = 10

			httpClient = createClient(maxConcurrentRequests, 100*time.Millisecond)

			// Make more requests that allowed in the window
			for range maxConcurrentRequests * 2 {
				go makeRequest(httpClient)
			}

			// Verify that only requests up to the limit was made immediately
			time.Sleep(5 * time.Millisecond)

			callCount := httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(maxConcurrentRequests))

			time.Sleep(100 * time.Millisecond)

			// verify that all requests was made
			callCount = httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(maxConcurrentRequests * 2))
		})
	})
})
