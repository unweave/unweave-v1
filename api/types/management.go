package types

import (
	"net/http"
	"regexp"
	"time"
)

const projectNameRegex = `^[\w.-]+$`

type AccessTokenCreateParams struct {
	Name string `json:"name"`
}

func (p *AccessTokenCreateParams) Bind(r *http.Request) error {
	if p.Name == "" {
		return &Error{
			Code:    http.StatusBadRequest,
			Message: "Name is required",
		}
	}
	return nil
}

type AccessTokenCreateResponse struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	Name      string    `json:"name"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type AccessTokensDeleteResponse struct {
	Success bool `json:"success"`
}

type Account struct {
	UserID              string    `json:"userID"`
	Email               string    `json:"email"`
	GithubID            int32     `json:"githubID"`
	GithubUsername      string    `json:"githubUsername"`
	DateJoined          time.Time `json:"dateJoined"`
	Credit              string    `json:"credit"`
	FirstName           string    `json:"firstName"`
	LastName            string    `json:"lastName"`
	Providers           []string  `json:"providers"`
	GithubCredentialsID *string   `json:"-,omitempty"`
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
	Name       string   `json:"name"`
	Tags       []string `json:"tags"`
	Visibility *string  `json:"visibility"`
	// SourceRepoCloneURL must be a HTTPS endpoint to a Git module i.e. https://github.com/unweave/unweave.git
	SourceRepoCloneURL *string `json:"repo"`
}

type ProjectCreateResponse struct {
	ID string `json:"id"`
}

func (p *ProjectCreateRequestParams) Bind(r *http.Request) error {
	if p.Name == "" {
		return &Error{
			Code:    http.StatusBadRequest,
			Message: "Name is required",
		}
	}

	regex := regexp.MustCompile(projectNameRegex)
	if !regex.MatchString(p.Name) {
		return &Error{
			Code:    http.StatusBadRequest,
			Message: "Name can only contain alphanumeric characters, underscores, dashes, and periods",
		}
	}

	if p.Visibility == nil {
		visibility := "private"
		p.Visibility = &visibility
	}
	return nil
}

type ProjectListResponse struct {
	Projects []Project `json:"projects"`
}

type ProjectGetResponse struct {
	Project Project `json:"project"`
}
