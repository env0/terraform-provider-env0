package http

//go:generate mockgen -destination=client_mock.go -package=http . HttpClientInterface

import (
	"github.com/go-resty/resty/v2"
)

type HttpClientInterface interface {
	Get(path string, params map[string]string, response interface{}) error
	Post(path string, request interface{}, response interface{}) error
	Put(path string, request interface{}, response interface{}) error
	Delete(path string) error
	Patch(path string, request interface{}, response interface{}) error
}

type HttpClient struct {
	jwtToken  string
	ApiKey    string
	ApiSecret string
	Endpoint  string
	client    *resty.Client
}

type HttpClientConfig struct {
	ApiKey      string
	ApiSecret   string
	ApiEndpoint string
	UserAgent   string
	RestClient  *resty.Client
}

func NewHttpClient(config HttpClientConfig) (*HttpClient, error) {
	httpClient := &HttpClient{
		ApiKey:    config.ApiKey,
		ApiSecret: config.ApiSecret,
		client:    config.RestClient.SetHostURL(config.ApiEndpoint).SetHeader("User-Agent", config.UserAgent),
	}
	var res string
	req := httpClient.client.R().SetBasicAuth(httpClient.ApiKey, httpClient.ApiSecret)
	response, err := req.SetQueryParams(map[string]string{"encoded": "true"}).SetResult(&res).Get("auth/token")
	if err != nil {
		return nil, err
	}

	httpClient.jwtToken = string(response.Body())
	return httpClient, nil
}

func (client *HttpClient) request() *resty.Request {
	return client.client.R().SetAuthToken(client.jwtToken)
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

func (client *HttpClient) Post(path string, request interface{}, response interface{}) error {
	result, err := client.request().
		SetBody(request).
		SetResult(response).
		Post(path)
	return client.httpResult(result, err)
}

func (client *HttpClient) Put(path string, request interface{}, response interface{}) error {
	result, err := client.request().
		SetBody(request).
		SetResult(response).
		Put(path)
	return client.httpResult(result, err)
}

func (client *HttpClient) Get(path string, params map[string]string, response interface{}) error {
	result, err := client.request().
		SetQueryParams(params).
		SetResult(response).
		Get(path)
	return client.httpResult(result, err)
}

func (client *HttpClient) Delete(path string) error {
	result, err := client.request().Delete(path)
	return client.httpResult(result, err)
}

func (client *HttpClient) Patch(path string, request interface{}, response interface{}) error {
	result, err := client.request().
		SetBody(request).
		SetResult(response).
		Patch(path)
	return client.httpResult(result, err)
}
