package types

import "net/http"

func (g *GithubIntegrationConnectRequest) Bind(r *http.Request) error {
	if g.Code == "" && g.AccessToken == "" {
		return &Error{
			Code:       http.StatusBadRequest,
			Message:    "Missing code or access token",
			Suggestion: "Either code or access token must be provided",
		}
	}

	if g.Code != "" && g.AccessToken != "" {
		return &Error{
			Code:       http.StatusBadRequest,
			Message:    "Both code and access token provided",
			Suggestion: "Either code or access token must be provided. Not both.",
		}
	}

	return nil
}
