package http_test

import (
	"net/http"
	"sync"
	"time"

	httpModule "github.com/env0/terraform-provider-env0/client/http"
	"github.com/env0/terraform-provider-env0/client/http/ratelimiter"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Keyed Rate Limiter", func() {
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

	createClientWithTotal := func(totalRequests, perPathRequests int, window time.Duration) *httpModule.HttpClient {
		config := httpModule.HttpClientConfig{
			ApiKey:      ApiKey,
			ApiSecret:   ApiSecret,
			ApiEndpoint: BaseUrl,
			UserAgent:   UserAgent,
			RestClient:  restClient,
			RateLimiter: ratelimiter.NewKeyedLimiter(totalRequests, window, func(string) int { return perPathRequests }),
		}
		client, err := httpModule.NewHttpClient(config)
		Expect(err).To(BeNil())

		return client
	}

	createClient := func(perPathRequests int, window time.Duration) *httpModule.HttpClient {
		return createClientWithTotal(1000, perPathRequests, window)
	}

	registerSuccess := func(method string, paths ...string) {
		for _, path := range paths {
			httpmock.RegisterResponder(method, BaseUrl+path, httpmock.NewStringResponder(200, SuccessResponse))
		}
	}

	// Sends the request in the background and returns a channel closed once it's answered.
	goRequest := func(method, path string) chan struct{} {
		done := make(chan struct{})

		go func() {
			defer close(done)

			var response string

			switch method {
			case http.MethodGet:
				_ = httpClient.Get(path, nil, &response)
			case http.MethodPost:
				_ = httpClient.Post(path, nil, &response)
			}
		}()

		return done
	}

	makeRequest := func(wg *sync.WaitGroup, client *httpModule.HttpClient) {
		defer GinkgoRecover()
		defer wg.Done()

		var response string

		err := client.Get(TestEndpoint, nil, &response)
		Expect(err).To(BeNil())
		Expect(response).To(Equal(SuccessResponse))
	}

	// Teardown resets the mock transport, so the spec has to outlive every request it started.
	requestsDone := func(wg *sync.WaitGroup) chan struct{} {
		done := make(chan struct{})

		go func() {
			wg.Wait()
			close(done)
		}()

		return done
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

			done := make(chan struct{})

			go func() {
				defer close(done)

				var response string

				_ = httpClient.Get(path, nil, &response)
			}()

			// Only the window's budget goes out immediately, the retries wait for a free slot.
			time.Sleep(20 * time.Millisecond)

			callCount := httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+path]).To(Equal(2))

			Eventually(done, 3*time.Second).Should(BeClosed())

			callCount = httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+path]).To(Equal(4))
		})

		// The WAF counts per method + path, so a 429 says the path is over its budget. Every request
		// to that path has to back off - retrying only the blocked one while its siblings keep firing
		// is what turns one 429 into thousands - but other paths have budget of their own.
		It("should pause only the path that answered 429", func() {
			const path = "/rate-limited"

			calls := 0

			httpmock.RegisterResponder("GET", BaseUrl+path,
				func(*http.Request) (*http.Response, error) {
					calls++
					if calls > 1 {
						return httpmock.NewStringResponse(200, SuccessResponse), nil
					}

					res := httpmock.NewStringResponse(http.StatusTooManyRequests, "TOO MANY REQUESTS")
					res.Header.Set("Retry-After", "1")

					return res, nil
				})

			httpClient = createClient(10, time.Minute)

			var rateLimitedResponse string

			Expect(httpClient.Get(path, nil, &rateLimitedResponse)).To(HaveOccurred())

			otherDone := goRequest(http.MethodGet, TestEndpoint)
			sameDone := goRequest(http.MethodGet, path)

			// A path that never answered 429 goes out right away, the 429'd one is held back.
			Eventually(otherDone, 100*time.Millisecond).Should(BeClosed())

			callCount := httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(1))
			Expect(callCount["GET "+BaseUrl+path]).To(Equal(1))

			// Retry-After said one second, after which the held back request goes out (plus the
			// limiter's wake-up spread).
			Eventually(sameDone, 3*time.Second).Should(BeClosed())

			callCount = httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+path]).To(Equal(2))
		})
	})

	Context("with rate limiting per method and path", func() {
		It("should not throttle requests to different paths against each other", func() {
			registerSuccess(http.MethodGet, "/a", "/b")

			httpClient = createClient(2, 200*time.Millisecond)

			firstDone, secondDone := goRequest(http.MethodGet, "/a"), goRequest(http.MethodGet, "/a")
			Eventually(firstDone, 100*time.Millisecond).Should(BeClosed())
			Eventually(secondDone, 100*time.Millisecond).Should(BeClosed())

			thirdDone := goRequest(http.MethodGet, "/a")
			otherDone := goRequest(http.MethodGet, "/b")

			Eventually(otherDone, 50*time.Millisecond).Should(BeClosed())
			Consistently(thirdDone, 50*time.Millisecond).ShouldNot(BeClosed())
			Eventually(thirdDone, 3*time.Second).Should(BeClosed())
		})

		// The WAF matches the URL path, which doesn't include the query string.
		It("should count requests to one path with different query strings against the same budget", func() {
			registerSuccess(http.MethodGet, "/a")

			httpClient = createClient(2, 200*time.Millisecond)

			firstDone, secondDone := goRequest(http.MethodGet, "/a?x=1"), goRequest(http.MethodGet, "/a?x=2")
			Eventually(firstDone, 100*time.Millisecond).Should(BeClosed())
			Eventually(secondDone, 100*time.Millisecond).Should(BeClosed())

			thirdDone := goRequest(http.MethodGet, "/a?x=3")

			Consistently(thirdDone, 50*time.Millisecond).ShouldNot(BeClosed())
			Eventually(thirdDone, 3*time.Second).Should(BeClosed())
		})

		It("should count different methods on the same path separately", func() {
			registerSuccess(http.MethodGet, "/a")
			registerSuccess(http.MethodPost, "/a")

			httpClient = createClient(1, time.Minute)

			getDone, postDone := goRequest(http.MethodGet, "/a"), goRequest(http.MethodPost, "/a")

			Eventually(getDone, 100*time.Millisecond).Should(BeClosed())
			Eventually(postDone, 100*time.Millisecond).Should(BeClosed())
		})

		It("should hold requests to every path to the total limit", func() {
			registerSuccess(http.MethodGet, "/a", "/b", "/c")

			httpClient = createClientWithTotal(2, 10, 200*time.Millisecond)

			aDone, bDone := goRequest(http.MethodGet, "/a"), goRequest(http.MethodGet, "/b")
			Eventually(aDone, 100*time.Millisecond).Should(BeClosed())
			Eventually(bDone, 100*time.Millisecond).Should(BeClosed())

			cDone := goRequest(http.MethodGet, "/c")

			Consistently(cDone, 50*time.Millisecond).ShouldNot(BeClosed())
			Eventually(cDone, 3*time.Second).Should(BeClosed())
		})
	})

	Context("with client rate limiting tests", func() {
		// These tests verify our HTTP client's rate limiting behavior
		It("should allow multiple requests up to the limit", func() {
			const maxConcurrentRequests = 10

			httpClient = createClient(maxConcurrentRequests, 100*time.Millisecond)

			var wg sync.WaitGroup

			wg.Add(maxConcurrentRequests)

			for range maxConcurrentRequests {
				go makeRequest(&wg, httpClient)
			}

			Eventually(requestsDone(&wg), 3*time.Second).Should(BeClosed())

			callCount := httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(maxConcurrentRequests))
		})

		It("should handle concurrent requests with rate limiting", func() {
			const maxConcurrentRequests = 10

			httpClient = createClient(maxConcurrentRequests, 100*time.Millisecond)

			var wg sync.WaitGroup

			wg.Add(maxConcurrentRequests * 2)

			// Make more requests that allowed in the window
			for range maxConcurrentRequests * 2 {
				go makeRequest(&wg, httpClient)
			}

			// Verify that only requests up to the limit was made immediately
			time.Sleep(5 * time.Millisecond)

			callCount := httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(maxConcurrentRequests))

			Eventually(requestsDone(&wg), 3*time.Second).Should(BeClosed())

			callCount = httpmock.GetCallCountInfo()
			Expect(callCount["GET "+BaseUrl+TestEndpoint]).To(Equal(maxConcurrentRequests * 2))
		})
	})
})
