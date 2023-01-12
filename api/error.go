package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/types"
)

type HTTPError struct {
	Code       int                   `json:"code"`
	Message    string                `json:"message"`
	Suggestion string                `json:"suggestion"`
	Provider   types.RuntimeProvider `json:"provider"`
	Err        error                 `json:"-"`
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *HTTPError) Render(w http.ResponseWriter, r *http.Request) error {
	// Depending on whether it is Unweave's fault or the user's fault, log the error
	// appropriately.
	if e.Code == http.StatusInternalServerError {
		log.Ctx(r.Context()).Error().Err(e).Stack().Msg(e.Message)
	} else {
		log.Ctx(r.Context()).Warn().Err(e).Stack().Msg(e.Message)
	}
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
			Err:        e.Err,
		}
	}
	return ErrInternalServer(err, fallbackMessage)
}

func ErrInternalServer(err error, msg string) render.Renderer {
	m := "Internal server error"
	if msg != "" {
		m = msg
	}
	return &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: m,
		Err:     err,
	}
}
