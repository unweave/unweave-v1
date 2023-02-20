package types

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

// Error
//
// Errors returned by the API should be as descriptive as possible and directly renderable
// to the consumer (CLI, web-app etc). Here are some examples:
//
// Provider errors
// ---------------
// Short:
//
//	LambdaLabs API error: Invalid Public Key
//
// Verbose:
//
//	LambdaLabs API error:
//		code: 400
//		message: Invalid Public Key
//	 	endpoint: POST /session
//
// Unweave errors
// --------------
// Short:
//
//	Unweave API error: Project not found
//
// Verbose:
//
//	Unweave API error:
//		code: 404
//		message: Project not found
//	 	endpoint: POST /session
//
// It should be possible to automatically generate the short and verbose versions of the
// error message from the same struct. The error message should not expose in inner workings
// of the API.
type Error struct {
	Code       int             `json:"code"`
	Message    string          `json:"message"`
	Suggestion string          `json:"suggestion,omitempty"`
	Provider   RuntimeProvider `json:"provider,omitempty"`
	Err        error           `json:"-"`
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *Error) Render(w http.ResponseWriter, r *http.Request) error {
	// Depending on whether it is Unweave's fault or the user's fault, log the error
	// appropriately.
	hook := log.Hook(NewErrLogHook())
	if e.Code == http.StatusInternalServerError {
		log.Ctx(r.Context()).Error().Err(e.Err).Stack().Msg(e.Message)
		hook.Error().Err(e.Err).Stack().Msg(e.Message)
	} else {
		log.Ctx(r.Context()).Warn().Err(e.Err).Stack().Msg(e.Message)
		hook.Warn().Err(e.Err).Stack().Msg(e.Message)
	}

	render.Status(r, e.Code)
	return nil
}

type UwError interface {
	Short() string
	Verbose() string
}
