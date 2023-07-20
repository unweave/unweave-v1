package types

import "net/http"

func (e *EndpointCreateParams) Bind(r *http.Request) error {
	if e.ExecID == "" {
		return &Error{
			Code:       http.StatusBadRequest,
			Message:    "Missing exec ID",
			Suggestion: "Exec ID must be provided",
		}
	}

	if e.Name == "" {
		return &Error{
			Code:       http.StatusBadRequest,
			Message:    "Missing name",
			Suggestion: "Name must be provided",
		}
	}

	return nil
}

func (e *EndpointVersionCreateParams) Bind(r *http.Request) error {
	if e.ExecID == "" {
		return &Error{
			Code:       http.StatusBadRequest,
			Message:    "Missing exec ID",
			Suggestion: "Exec ID must be provided",
		}
	}

	return nil
}
