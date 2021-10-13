package client

import "net/url"

type parameter struct {
	key   string
	value string
}

func newQueryURL(path string, params ...parameter) (*url.URL, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	q := u.Query()

	for _, param := range params {
		q.Add(param.key, param.value)
	}

	u.RawQuery = q.Encode()

	return u, nil
}
