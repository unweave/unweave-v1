package types

import (
	"github.com/google/go-github/v53/github"
)

// CloneURL represents the URL Git repositories should be used to clone to
type CloneURL = string

type Repository struct {
	Name     string   `json:"name,omitempty"`
	CloneURL CloneURL `json:"cloneURL,omitempty"`
}

func NewRepository(gh *github.Repository) Repository {
	return Repository{
		Name:     gh.GetName(),
		CloneURL: gh.GetCloneURL(),
	}
}

func NewRepositories(gh []*github.Repository) []Repository {
	out := make([]Repository, 0, len(gh))

	for _, repository := range gh {
		if repository != nil {
			out = append(out, NewRepository(repository))
		}
	}

	return out
}

// GithubListRepositoriesResponse lists all repositories authenticated to Unweave,
// whether the GitHub application is installed,
// and the installation URL if applicable.
type GithubListRepositoriesResponse struct {
	Repositories []Repository `json:"repositories,omitempty"`
	Installed    bool         `json:"installed,omitempty"`
	InstallURL   string       `json:"installURL,omitempty"`
}
