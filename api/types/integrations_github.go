package types

// CloneURL represents the URL Git repositories should be used to clone to
type CloneURL = string

type Repository struct {
	Name     string   `json:"name,omitempty"`
	CloneURL CloneURL `json:"cloneURL,omitempty"`
}

// GithubListRepositoriesResponse lists all repositories authenticated to Unweave,
// whether the GitHub application is installed,
// and the installation URL if applicable.
type GithubListRepositoriesResponse struct {
	Repositories []Repository `json:"repositories,omitempty"`
	Installed    bool         `json:"installed,omitempty"`
	InstallURL   string       `json:"installURL,omitempty"`
}
