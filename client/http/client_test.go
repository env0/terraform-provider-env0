package http_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	httpModule "github.com/env0/terraform-provider-env0/client/http"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

// The JSON here returns camelCase keys
type ResponseType struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type RequestBody struct {
	Message string `json:"message"`
}

const BaseUrl = "https://fake.env0.com"
const ApiKey = "MY_USER"
const ApiSecret = "MY_PASS"
const ExpectedBasicAuth = "Basic TVlfVVNFUjpNWV9QQVNT"
const UserAgent = "super-cool-ua"
const ErrorStatusCode = 500
const ErrorMessage = "Very bad!"

var httpclient *httpModule.HttpClient

var _ = BeforeSuite(func() {
	// mock all HTTP requests
	restClient := resty.New()
	config := httpModule.HttpClientConfig{
		ApiKey:      ApiKey,
		ApiSecret:   ApiSecret,
		ApiEndpoint: BaseUrl,
		UserAgent:   UserAgent,
		RestClient:  restClient,
	}
	httpclient, _ = httpModule.NewHttpClient(config)
	httpmock.ActivateNonDefault(restClient.GetClient())
})

var _ = BeforeEach(func() {
	httpmock.Reset()
})

var _ = AfterSuite(func() {
	// unmock HTTP requests
	httpmock.DeactivateAndReset()
})

func TestHttpClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTP Client Tests")
}

var _ = Describe("Http Client", func() {
	var httpRequest *http.Request

	mockRequest := RequestBody{
		Message: "I have a request",
	}
	mockedResponse := ResponseType{
		Id:   1,
		Name: "I have a response",
	}
	successURI := "/path/to/success"
	failureURI := "/path/to/failure"
	successUrl := BaseUrl + successURI
	failureUrl := BaseUrl + failureURI

	AssertAuth := func() {
		authorization := httpRequest.Header["Authorization"]
		Expect(len(authorization)).To(Equal(1), "Should have authorization header")
		Expect(authorization[0]).To(Equal(ExpectedBasicAuth), "Should have correct Basic Auth")
	}

	AssertNoError := func(err error) {
		Expect(err).To(BeNil(), "Should not return error")
	}

	AssertError := func(err error) {
		Expect(err.Error()).To(Equal(strconv.Itoa(ErrorStatusCode)+": "+ErrorMessage), "Should return error message")
	}

	AssertHttpCall := func(method string, url string) {
		methodAndUrl := method + " " + url
		// Validate call happened once
		callMap := httpmock.GetCallCountInfo()
		Expect(callMap[methodAndUrl]).Should(Equal(1), "Should call "+methodAndUrl)

		// Validate no other call happened
		delete(callMap, methodAndUrl)
		for unexpectedCall, amount := range callMap {
			Expect(amount).To(BeZero(), "Should not call "+unexpectedCall)
		}
	}

	AssertReturnedResponse := func(result ResponseType) {
		Expect(result).To(Equal(mockedResponse), "Should return expected response")
	}

	AssertRequestBody := func() {
		// Convert mock request to JSON
		mockRequestJson, _ := json.Marshal(mockRequest)

		// Get actual request body
		actualBodyBuffer := new(strings.Builder)
		_, _ = io.Copy(actualBodyBuffer, httpRequest.Body)

		Expect(actualBodyBuffer.String()).To(Equal(string(mockRequestJson)), "Should send payload as HTTP request body")
	}

	BeforeEach(func() {
		httpRequest = nil
		// Make calls to /path/to/success return 200, and calls to /path/to/failure return 500
		for _, methodType := range []string{"GET", "POST", "PUT", "DELETE"} {
			httpmock.RegisterResponder(methodType, successUrl, func(req *http.Request) (*http.Response, error) {
				httpRequest = req

				return httpmock.NewJsonResponse(200, mockedResponse)
			})
			httpmock.RegisterResponder(methodType, failureUrl, func(req *http.Request) (*http.Response, error) {
				httpRequest = req

				return httpmock.NewStringResponse(ErrorStatusCode, ErrorMessage), nil
			})
		}
	})

	Describe("Get", func() {
		DescribeTable("2XX response",
			func(params map[string]string, expectedQuery string) {
				var result ResponseType
				var err = httpclient.Get(successURI, params, &result)

				AssertHttpCall("GET", successUrl)
				AssertNoError(err)

				AssertReturnedResponse(result)
				AssertAuth()
				Expect(httpRequest.URL.RawQuery).To(Equal(expectedQuery), "Should parse query params")
			},
			Entry("Without params", nil, ""),
			Entry("With params", map[string]string{
				"param1": "1",
				"param2": "two",
			}, "param1=1&param2=two"),
		)

		It("5XX response", func() {
			var result ResponseType
			var err = httpclient.Get(failureURI, nil, &result)

			AssertHttpCall("GET", failureUrl)
			AssertError(err)
		})
	})

	Describe("Post", func() {
		It("2XX response", func() {
			var result ResponseType
			var err = httpclient.Post(successURI, mockRequest, &result)

			AssertHttpCall("POST", successUrl)
			AssertNoError(err)

			AssertReturnedResponse(result)
			AssertRequestBody()
			AssertAuth()
		})

		It("5XX response", func() {
			var result ResponseType
			var err = httpclient.Post(failureURI, mockRequest, &result)

			AssertHttpCall("POST", failureUrl)
			AssertError(err)
		})
	})

	Describe("Delete", func() {
		It("2XX response", func() {
			var err = httpclient.Delete(successURI, nil)

			AssertHttpCall("DELETE", successUrl)
			AssertNoError(err)
			AssertAuth()
		})

		It("5XX response", func() {
			var err = httpclient.Delete(failureURI, nil)

			AssertHttpCall("DELETE", failureUrl)
			AssertError(err)
		})
	})

	Describe("Put", func() {
		It("2XX response", func() {
			var result ResponseType
			var err = httpclient.Put(successURI, mockRequest, &result)

			AssertHttpCall("PUT", successUrl)
			AssertNoError(err)

			AssertReturnedResponse(result)
			AssertRequestBody()
			AssertAuth()
		})

		It("5XX response", func() {
			var result ResponseType
			var err = httpclient.Put(failureURI, mockRequest, &result)

			AssertHttpCall("PUT", failureUrl)
			AssertError(err)
		})
	})
})
