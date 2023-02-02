package types

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

type AccessTokenCreateParams struct {
	Name string `json:"name"`
}

func (p *AccessTokenCreateParams) Bind(r *http.Request) error {
	if p.Name == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Name is required",
		}
	}
	return nil
}

type AccessTokenCreateResponse struct {
	ID        uuid.UUID `json:"id"`
	Token     string    `json:"token"`
	Name      string    `json:"name"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type AccessTokensDeleteResponse struct {
	Success bool `json:"success"`
}

type Account struct {
	UserID         uuid.UUID `json:"userID"`
	Email          string    `json:"email"`
	GithubID       int32     `json:"githubID"`
	GithubUsername string    `json:"githubUsername"`
	DateJoined     time.Time `json:"dateJoined"`
	Credit         string    `json:"credit"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Providers      []string  `json:"providers"`
}

type AccountGetResponse struct {
	Account Account `json:"account"`
}

type PairingTokenCreateResponse struct {
	Code string `json:"code"`
}

type PairingTokenExchangeResponse struct {
	Token   string  `json:"token"`
	Account Account `json:"account"`
}

type ProjectCreateRequestParams struct {
	Name string `json:"name"`
}

type ProjectCreateResponse struct {
	ID uuid.UUID `json:"id"`
}

func (p *ProjectCreateRequestParams) Bind(r *http.Request) error {
	if p.Name == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Name is required",
		}
	}
	return nil
}

type ProjectListResponse struct {
	Projects []Project `json:"projects"`
}

type ProjectGetResponse struct {
	Project Project `json:"project"`
}
