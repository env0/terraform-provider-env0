package http

import (
	"net/http"

	"github.com/go-resty/resty/v2"
)

func NewMockFailedResponseError(statusCode int) error {
	raw := &http.Response{
		StatusCode: statusCode,
	}

	res := &resty.Response{
		RawResponse: raw,
	}

	return &FailedResponseError{
		res: res,
	}
}
