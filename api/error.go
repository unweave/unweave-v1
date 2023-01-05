package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/unweave/unweave/types"
)

type HTTPError struct {
	Code       int                   `json:"code"`
	Message    string                `json:"message"`
	Suggestion string                `json:"suggestion"`
	Provider   types.RuntimeProvider `json:"provider"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Code)
	return nil
}

func ErrHTTPError(err error, fallbackMessage string) render.Renderer {
	if err == nil {
		return nil
	}
	var e *types.Error
	if errors.As(err, &e) {
		return &HTTPError{
			Code:       e.Code,
			Message:    e.Message,
			Provider:   e.Provider,
			Suggestion: e.Suggestion,
		}
	}
	return ErrInternalServer(fallbackMessage)
}

func ErrInternalServer(msg string) render.Renderer {
	m := "Internal server error"
	if msg != "" {
		m = msg
	}
	return &HTTPError{
		Code:    500,
		Message: m,
	}
}
