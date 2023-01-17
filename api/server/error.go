package server

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/unweave/unweave/api/types"
)

func ErrHTTPError(err error, fallbackMessage string) render.Renderer {
	if err == nil {
		return nil
	}
	var e *types.Error
	if errors.As(err, &e) {
		return &types.HTTPError{
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
	return &types.HTTPError{
		Code:    http.StatusInternalServerError,
		Message: m,
		Err:     err,
	}
}
