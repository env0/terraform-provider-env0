package http

import "github.com/go-resty/resty/v2"

type FailedResponseError struct {
	res *resty.Response
}

func (e *FailedResponseError) Error() string {
	return e.res.Status() + ": " + string(e.res.Body())
}

func (e *FailedResponseError) NotFound() bool {
	return e.res.StatusCode() == 404
}

func (e *FailedResponseError) BadRequest() bool {
	return e.res.StatusCode() == 400
}
