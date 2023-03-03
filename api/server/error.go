package server

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/unweave/unweave/api/types"
)

func ErrHTTPBadRequest(err error, fallbackMessage string) render.Renderer {
	var e *types.Error
	if errors.As(err, &e) {
		return &types.Error{
			Code:       e.Code,
			Message:    e.Message,
			Provider:   e.Provider,
			Suggestion: e.Suggestion,
			Err:        err,
		}
	}
	return &types.Error{
		Code:    http.StatusBadRequest,
		Message: fallbackMessage,
		Err:     err,
	}
}

func ErrHTTPError(err error, fallbackMessage string) render.Renderer {
	if err == nil {
		return nil
	}
	var e *types.Error
	if errors.As(err, &e) {
		return &types.Error{
			Code:       e.Code,
			Message:    e.Message,
			Provider:   e.Provider,
			Suggestion: e.Suggestion,
			Err:        err,
		}
	}
	return ErrInternalServer(err, fallbackMessage)
}

func ErrInternalServer(err error, msg string) render.Renderer {
	m := "Internal server error"
	if msg != "" {
		m = msg
	}
	return &types.Error{
		Code:    http.StatusInternalServerError,
		Message: m,
		Err:     err,
	}
}
