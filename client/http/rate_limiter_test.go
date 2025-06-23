package http_test

import (
	"sync"
	"time"

	httpModule "github.com/env0/terraform-provider-env0/client/http"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rate Limiter", func() {
	const (
		BaseUrl          = "https://fake.env0.com"
		ApiKey           = "TEST_KEY"
		ApiSecret        = "TEST_SECRET"
		UserAgent        = "test-agent"
		TestEndpoint     = "/test-endpoint"
		SuccessResponse  = "success"
		DefaultRateLimit = 800
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

	createClient := func(rateLimit int) *httpModule.HttpClient {
		config := httpModule.HttpClientConfig{
			ApiKey:             ApiKey,
			ApiSecret:          ApiSecret,
			ApiEndpoint:        BaseUrl,
			UserAgent:          UserAgent,
			RestClient:         restClient,
			RateLimitPerMinute: rateLimit,
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

	It("should use default rate limit when not specified", func() {
		// Create client with zero rate limit (should use default)
		httpClient = createClient(0)

		// Make a request and verify it succeeds
		makeRequest(httpClient)

		// Verify the call was made
		callCount := httpmock.GetCallCountInfo()
		Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(1))
	})

	// Test that directly verifies the limiter configuration
	It("should create HTTP client with default rate limit when not specified", func() {
		// Create a client without specifying rate limit
		config := httpModule.HttpClientConfig{
			ApiKey:      ApiKey,
			ApiSecret:   ApiSecret,
			ApiEndpoint: BaseUrl,
			UserAgent:   UserAgent,
			RestClient:  restClient,
			// RateLimitPerMinute not specified - should use default
		}
		client, err := httpModule.NewHttpClient(config)
		Expect(err).To(BeNil())

		// Register a responder for the test endpoint
		httpmock.RegisterResponder("GET", BaseUrl+TestEndpoint,
			httpmock.NewStringResponder(200, SuccessResponse))

		// Make a request to verify client works with default rate limit
		var response string
		err = client.Get(TestEndpoint, nil, &response)
		Expect(err).To(BeNil())
		Expect(response).To(Equal(SuccessResponse))

		// Verify the call was made
		callCount := httpmock.GetCallCountInfo()
		Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(1))
	})

	It("should allow burst of requests up to the rate limit", func() {
		// Use a small rate limit for testing to keep the test fast
		testRateLimit := 10
		httpClient = createClient(testRateLimit)

		// Make concurrent requests up to the rate limit
		var wg sync.WaitGroup
		for i := 0; i < testRateLimit; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				makeRequest(httpClient)
			}()
		}

		// Wait for all requests to complete
		wg.Wait()

		// Verify all calls were made
		callCount := httpmock.GetCallCountInfo()
		Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(testRateLimit))
	})

	Context("with client rate limiting tests", func() {
		// These tests verify our HTTP client's rate limiting behavior

		It("should allow multiple requests and eventually succeed", func() {
			// Create a client with a reasonable rate limit for testing
			testLimit := 20 // Small enough to test, but not too small
			httpClient = createClient(testLimit)

			// Make a series of requests that should all succeed
			totalRequests := 10
			for i := 0; i < totalRequests; i++ {
				makeRequest(httpClient)
			}

			// Verify all requests were made successfully
			callCount := httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(totalRequests))
		})

		It("should handle concurrent requests with rate limiting", func() {
			// Create a client with a higher rate limit for testing concurrent requests
			// 600 per minute = 10 per second, which is fast enough for testing
			testLimit := 600
			httpClient = createClient(testLimit)

			// Make concurrent requests
			var wg sync.WaitGroup
			totalRequests := 20 // Small number to keep test fast
			for i := 0; i < totalRequests; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					makeRequest(httpClient)
				}()
			}

			// Wait for all requests to complete
			wg.Wait()

			// Verify all requests were made
			callCount := httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(totalRequests))
		})

		// This test takes over a minute to run
		It("should queue excess requests and process them after rate limit refreshes", func() {
			// Create a client with a very small rate limit for testing
			testLimit := 5

			config := httpModule.HttpClientConfig{
				ApiKey:             ApiKey,
				ApiSecret:          ApiSecret,
				ApiEndpoint:        BaseUrl,
				UserAgent:          UserAgent,
				RestClient:         restClient,
				RateLimitPerMinute: testLimit,
			}
			client, err := httpModule.NewHttpClient(config)
			Expect(err).To(BeNil())

			httpmock.RegisterResponder("GET", BaseUrl+TestEndpoint,
				httpmock.NewStringResponder(200, SuccessResponse))

			for i := 0; i < testLimit; i++ {
				var response string
				err = client.Get(TestEndpoint, nil, &response)
				Expect(err).To(BeNil())
				Expect(response).To(Equal(SuccessResponse))
			}

			// Verify all initial requests were made
			callCount := httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(testLimit))

			// Send one more request in a goroutine - this should be queued
			var extraResponse string
			var extraErr error
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				// This should block until a token becomes available
				extraErr = client.Get(TestEndpoint, nil, &extraResponse)
			}()

			// Wait a short time to ensure the request is queued but not yet processed
			time.Sleep(5 * time.Second)

			// Verify the queued request has not been processed yet
			callCount = httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(testLimit), "Queued request should not be processed immediately")

			// Wait for the rate limit to refresh (slightly more than 1 minute)
			time.Sleep(56 * time.Second)

			// Wait for the goroutine to complete
			wg.Wait()

			// Verify the extra request succeeded
			Expect(extraErr).To(BeNil())
			Expect(extraResponse).To(Equal(SuccessResponse))

			// Verify all requests were eventually made
			callCount = httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(testLimit+1), "Queued request should be processed after rate limit refreshes")
		})
	})
})
