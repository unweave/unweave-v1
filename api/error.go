package api

import (
	"net/http"

	"github.com/go-chi/render"
)

// swagger:response errorResponse
type httpError struct {
	// in: body
	Body HTTPError
}

type HTTPError struct {
	// Example: 400
	Code int `json:"code"`
	// Example: Bad Request
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

func ErrInternalServer(msg string) render.Renderer {
	return &HTTPError{
		Code:    500,
		Message: msg,
	}
}
