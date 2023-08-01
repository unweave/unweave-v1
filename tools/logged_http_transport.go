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
	ctx := log.Ctx(req.Context()).With().
		Str("method", req.Method).
		Str("host", req.Host).
		Str("path", req.URL.Path).
		Logger().WithContext(req.Context())

	log.Ctx(ctx).Debug().
		Str("event", "http_request").
		Send()

	startTime := time.Now()
	resp, err := t.base.RoundTrip(req)
	duration := time.Since(startTime)

	event := log.Ctx(ctx).Debug().
		Str("event", "http_response").
		Dur("duration", duration).
		Err(err)

	if err == nil {
		event = event.Int("status", resp.StatusCode)
	}

	event.Send()

	return resp, err
}
