package http

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
)

// The WAF counts requests per HTTP method + IP + URL path (env0 prerequisites/resources/waf.yml),
// so the client keys its limiter the same way. The query string isn't part of the path.
//
// The before-request hook sees the path relative to the base URL, the after-response hook sees the
// absolute URL, and both have to land on the same key: a 429 pause on a key Wait never checks
// wouldn't hold anything back.
func rateLimitKey(c *resty.Client, method, rawURL string) string {
	u, err := url.Parse(rawURL)
	if err == nil && !u.IsAbs() {
		// Joined the same way resty resolves a relative request URL.
		relative := u.String()
		if !strings.HasPrefix(relative, "/") {
			relative = "/" + relative
		}

		u, err = url.Parse(c.BaseURL + relative)
	}

	if err != nil {
		return method + " " + rawURL
	}

	return method + " " + u.EscapedPath()
}

// PerPathRequestLimit returns the per window limit for a rate limiter key: criticalPath for the
// paths the WAF's critical-path rule covers, perPath for every other one.
func PerPathRequestLimit(perPath, criticalPath int) func(key string) int {
	return func(key string) int {
		method, path, _ := strings.Cut(key, " ")
		if isCriticalPath(method, path) {
			return criticalPath
		}

		return perPath
	}
}

// Mirrors the match statements of the WAF's CriticalPathRateLimit rule, which counts these paths on
// top of the regular per path rule with a much lower limit.
func isCriticalPath(method, path string) bool {
	switch method {
	case http.MethodGet:
		return path == "/environments"
	case http.MethodPost:
		return strings.HasPrefix(path, "/environments")
	case http.MethodPut:
		return strings.Contains(path, "/environments/deployments/")
	}

	return false
}
