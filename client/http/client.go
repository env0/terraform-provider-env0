package http

//go:generate mockgen -destination=client_mock.go -package=http . HttpClientInterface

import (
	"context"
	"reflect"

	"github.com/env0/terraform-provider-env0/client/http/ratelimiter"
	"github.com/go-resty/resty/v2"
)

type HttpClientInterface interface {
	Get(path string, params map[string]string, response any) error
	Post(path string, request any, response any) error
	Put(path string, request any, response any) error
	Delete(path string, params map[string]string) error
	Patch(path string, request any, response any) error
}

type HttpClient struct {
	ApiKey      string
	ApiSecret   string
	Endpoint    string
	client      *resty.Client
	rateLimiter *ratelimiter.RateLimiter
}

type HttpClientConfig struct {
	ApiKey      string
	ApiSecret   string
	ApiEndpoint string
	UserAgent   string
	RestClient  *resty.Client
	RateLimiter ratelimiter.RateLimiter // Optional rate limiter to limit outgoing requests
}

func NewHttpClient(config HttpClientConfig) (*HttpClient, error) {
	httpClient := &HttpClient{
		ApiKey:      config.ApiKey,
		ApiSecret:   config.ApiSecret,
		client:      config.RestClient.SetBaseURL(config.ApiEndpoint).SetHeader("User-Agent", config.UserAgent),
		rateLimiter: &config.RateLimiter,
	}

	return httpClient, nil
}

func (client *HttpClient) request() *resty.Request {
	if *client.rateLimiter != nil {
		ctx := context.Background()

		err := (*client.rateLimiter).Wait(ctx)
		if err != nil {
			return client.client.R().SetError(err)
		}
	}

	return client.client.R().SetBasicAuth(client.ApiKey, client.ApiSecret)
}

func (client *HttpClient) httpResult(response *resty.Response, err error) error {
	if err != nil {
		return err
	}

	if !response.IsSuccess() {
		return &FailedResponseError{res: response}
	}

	return nil
}

func (client *HttpClient) Post(path string, request any, response any) error {
	req := client.request().SetBody(request)
	if response != nil {
		req = req.SetResult(response)
	}

	result, err := req.Post(path)

	return client.httpResult(result, err)
}

func (client *HttpClient) Put(path string, request any, response any) error {
	req := client.request().SetBody(request)
	if response != nil {
		req = req.SetResult(response)
	}

	result, err := req.Put(path)

	return client.httpResult(result, err)
}

func (client *HttpClient) Get(path string, params map[string]string, response any) error {
	request := client.request().SetQueryParams(params)

	responseType := reflect.TypeOf(response)

	if responseType.Kind() == reflect.Ptr && responseType.Elem().Kind() == reflect.String {
		responseStrPtr := response.(*string)

		result, err := request.Get(path)
		if err == nil {
			*responseStrPtr = string(result.Body())
		}

		return client.httpResult(result, err)
	} else {
		return client.httpResult(request.SetResult(response).Get(path))
	}
}

func (client *HttpClient) Delete(path string, params map[string]string) error {
	result, err := client.request().SetQueryParams(params).Delete(path)

	return client.httpResult(result, err)
}

func (client *HttpClient) Patch(path string, request any, response any) error {
	result, err := client.request().
		SetBody(request).
		SetResult(response).
		Patch(path)

	return client.httpResult(result, err)
}
