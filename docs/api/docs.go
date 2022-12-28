// Package classification Unweave API.
//
// Documentation of the Unweave API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//	Host: api.unweave.io
//	License: Apache 2.0 https://www.apache.org/licenses/LICENSE-2.0.html
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package docs

import "github.com/unweave/unweave-v2/api"

// swagger:route POST /session session create
// responses:
// 	201: sessionCreate

// swagger:parameters session create
type sessionCreateRequest struct {
	// in: body
	Body api.SessionCreateParams
}

// swagger:response sessionCreate
type sessionCreateResponse struct {
	// in: body
	Body api.Session
}

// swagger:route GET /session/{id} session get-session
// parameters:
// + name: id
//   in: path
//   required: true
//   type: string
//   description: session ID
// responses:
// 	200: sessionGetResponse

// swagger:response sessionGetResponse
type sessionGetResponse struct {
	// in: body
	Body api.Session
}

// swagger:response errorResponse
type httpError struct {
	// in: body
	Body api.HTTPError
}
