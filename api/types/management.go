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

type PairingTokenCreateResponse struct {
	Code string `json:"code"`
}

type PairingTokenExchangeResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
}
