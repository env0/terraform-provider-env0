package http

import (
	"net/http"
	"strconv"
	"time"
)

// ParseRetryAfter parses the value of a Retry-After response header, which is either a number
// of seconds or an HTTP-date (RFC 7231 section 7.1.3). The second return value is false when
// the header is missing, unparsable or already in the past, in which case the caller should
// fall back to its own backoff.
func ParseRetryAfter(header string) (time.Duration, bool) {
	if header == "" {
		return 0, false
	}

	if seconds, err := strconv.Atoi(header); err == nil {
		if seconds <= 0 {
			return 0, false
		}

		return time.Duration(seconds) * time.Second, true
	}

	date, err := http.ParseTime(header)
	if err != nil {
		return 0, false
	}

	delay := time.Until(date)
	if delay <= 0 {
		return 0, false
	}

	return delay, true
}
