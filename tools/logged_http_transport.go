package tools

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type transport struct {
	base http.RoundTripper
}

func LoggedHTTPTransport(t http.RoundTripper) http.RoundTripper {
	return &transport{base: t}
}

//nolint:wrapcheck
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Debug().
		Str("method", req.Method).
		Str("host", req.Host).
		Str("path", req.URL.Path).
		Str("event", "http_request").
		Send()

	startTime := time.Now()
	resp, err := t.base.RoundTrip(req)
	duration := time.Since(startTime)

	entry := log.Debug().
		Str("method", req.Method).
		Str("host", req.Host).
		Str("path", req.URL.Path).
		Str("event", "http_response").
		Dur("duration", duration).
		Err(err)
	if err == nil {
		entry.Int("status", resp.StatusCode)
	}

	entry.Send()

	return resp, err
}
