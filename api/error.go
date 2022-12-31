package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/unweave/unweave-v2/types"
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

func (e *HTTPError) Short() string {
	str := fmt.Sprintf("%s API error: %s", e.Provider, e.Message)
	return str
}

func (e *HTTPError) Verbose() string {
	str := "API error:\n"
	if e.Provider != "" {
		str = fmt.Sprintf("%s API error:\n", e.Provider)
	}
	if e.Code != 0 {
		str += fmt.Sprintf("  Code: %d\n", e.Code)
	}
	if e.Message != "" {
		str += fmt.Sprintf("  Message: %s\n", e.Message)
	}
	if e.Suggestion != "" {
		str += fmt.Sprintf("  Suggestion: %s\n", e.Suggestion)
	}
	return str
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

func ErrBadRequest(msg string) render.Renderer {
	return &HTTPError{
		Code:    400,
		Message: msg,
	}
}

func ErrInternalServer(msg string) render.Renderer {
	return &HTTPError{
		Code:    500,
		Message: msg,
	}
}
