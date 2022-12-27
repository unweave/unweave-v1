package api

import (
	"net/http"

	"github.com/go-chi/render"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *HTTPError) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Code)
	return nil
}

func ErrBadRequest(msg string) render.Renderer {
	return &HTTPError{
		Code:    400,
		Message: msg,
	}
}
