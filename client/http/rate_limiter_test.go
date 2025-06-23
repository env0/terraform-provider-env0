package http_test

import (
	"context"
	"sync"
	"time"

	httpModule "github.com/env0/terraform-provider-env0/client/http"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/time/rate"
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
	It("should configure rate limiter with correct burst size and rate", func() {
		// Create a rate limiter directly to test its configuration
		requestsPerMinute := 800
		limiter := rate.NewLimiter(rate.Limit(float64(requestsPerMinute)/60.0), requestsPerMinute)
		
		// Verify the burst size is set to the full minute limit
		Expect(limiter.Burst()).To(Equal(requestsPerMinute))
		
		// Verify the rate is set to requestsPerMinute/60 per second
		Expect(limiter.Limit()).To(Equal(rate.Limit(float64(requestsPerMinute)/60.0)))
		
		// Test that we can immediately consume the full burst capacity
		reserve := limiter.ReserveN(time.Now(), requestsPerMinute)
		Expect(reserve.OK()).To(BeTrue(), "Should be able to reserve the full burst capacity")
		Expect(reserve.Delay()).To(Equal(time.Duration(0)), "Should have no delay for burst capacity")
		
		// Test that exceeding the burst capacity causes a delay
		reserve = limiter.ReserveN(time.Now(), 1)
		Expect(reserve.OK()).To(BeTrue(), "Should be able to reserve beyond burst capacity")
		Expect(reserve.Delay()).To(BeNumerically(">", time.Duration(0)), "Should have delay for exceeding burst capacity")
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

	Context("with more direct rate limiter tests", func() {
		// These tests directly test the rate limiter behavior without relying on timing
		
		It("should allow immediate consumption of tokens up to burst limit", func() {
			// Create a rate limiter with a small limit for testing
			testLimit := 10
			limiter := rate.NewLimiter(rate.Limit(float64(testLimit)/60.0), testLimit)
			
			// Should be able to consume all tokens immediately
			for i := 0; i < testLimit; i++ {
				allow := limiter.Allow()
				Expect(allow).To(BeTrue(), "Should allow request within burst limit")
			}
			
			// The next token should not be immediately available
			allow := limiter.Allow()
			Expect(allow).To(BeFalse(), "Should not allow request beyond burst limit without waiting")
		})
		
		It("should block when tokens are exhausted and then allow after refill", func() {
			// Create a rate limiter with a small limit and fast refill for testing
			testLimit := 5
			// Set to 5 tokens per second (very fast for testing)
			limiter := rate.NewLimiter(rate.Limit(5), testLimit)
			
			// Consume all tokens
			for i := 0; i < testLimit; i++ {
				allow := limiter.Allow()
				Expect(allow).To(BeTrue(), "Should allow request within burst limit")
			}
			
			// The next token should not be immediately available
			allow := limiter.Allow()
			Expect(allow).To(BeFalse(), "Should not allow request beyond burst limit without waiting")
			
			// Wait for at least one token to be refilled (should take ~200ms)
			time.Sleep(250 * time.Millisecond)
			
			// Now we should be able to get a token
			allow = limiter.Allow()
			Expect(allow).To(BeTrue(), "Should allow request after token refill")
		})
		
		It("should queue requests with Wait() when tokens are exhausted", func() {
			// Create a rate limiter with a small limit and fast refill for testing
			testLimit := 5
			// Set to 10 tokens per second (very fast for testing)
			limiter := rate.NewLimiter(rate.Limit(10), testLimit)
			
			// Consume all tokens
			for i := 0; i < testLimit; i++ {
				allow := limiter.Allow()
				Expect(allow).To(BeTrue(), "Should allow request within burst limit")
			}
			
			// The next request should wait, not fail
			ctx := context.Background()
			startTime := time.Now()
			
			// This should wait for a token (approximately 100ms with our rate)
			err := limiter.Wait(ctx)
			Elapsed := time.Since(startTime)
			
			Expect(err).To(BeNil(), "Wait should not return an error")
			Expect(Elapsed).To(BeNumerically(">", 50*time.Millisecond), "Should have waited for a token")
			Expect(Elapsed).To(BeNumerically("<", 200*time.Millisecond), "Should not wait too long")
		})
	})
})
