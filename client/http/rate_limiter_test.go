package http_test

import (
	"net/http"
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

	Context("with rate limiting applied to retries", func() {
		// Retried attempts used to bypass the limiter, so a run that started failing sent its
		// retries on top of the allowed rate - which is what kept the WAF blocking requests.
		It("should count retried attempts against the limit", func() {
			const path = "/server-error"

			httpmock.RegisterResponder("GET", BaseUrl+path,
				httpmock.NewStringResponder(http.StatusInternalServerError, "BAD"))

			restClient.SetRetryCount(3).
				SetRetryWaitTime(time.Millisecond).
				SetRetryMaxWaitTime(time.Millisecond).
				AddRetryCondition(func(r *resty.Response, err error) bool {
					return r.StatusCode() >= http.StatusInternalServerError
				})

			httpClient = createClient(2, 200*time.Millisecond)

			var response string

			go func() {
				_ = httpClient.Get(path, nil, &response)
			}()

			// Only the window's budget goes out immediately, the retries wait for a free slot.
			time.Sleep(20 * time.Millisecond)

			callCount := httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+path]).To(Equal(2))

			time.Sleep(300 * time.Millisecond)

			callCount = httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+path]).To(Equal(4))
		})

		// A 429 is a server side limit, so the entire client has to back off. Retrying only the
		// blocked request while its siblings keep firing is what turns one 429 into thousands.
		It("should pause every request after a 429", func() {
			const path = "/rate-limited"

			httpmock.RegisterResponder("GET", BaseUrl+path,
				func(*http.Request) (*http.Response, error) {
					res := httpmock.NewStringResponse(http.StatusTooManyRequests, "TOO MANY REQUESTS")
					res.Header.Set("Retry-After", "1")

					return res, nil
				})

			httpClient = createClient(10, time.Minute)

			var response string

			Expect(httpClient.Get(path, nil, &response)).ToNot(BeNil())

			go func() {
				_ = httpClient.Get(TestEndpoint, nil, &response)
			}()

			time.Sleep(100 * time.Millisecond)

			callCount := httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(0))

			// Retry-After said one second, after which the held back request goes out.
			time.Sleep(time.Second)

			callCount = httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(1))
		})
	})

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
