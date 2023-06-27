package types

type Repository struct {
	Name string `json:"name,omitempty"`
	// URL must be cloneable by Git
	URL URL `json:"url,omitempty"`
}

// GithubListRepositoriesResponse lists all repositories authenticated to Unweave,
// whether the GitHub application is installed,
// and the installation URL if applicable.
type GithubListRepositoriesResponse struct {
	Repositories []Repository `json:"repositories,omitempty"`
	Installed    bool         `json:"installed,omitempty"`
	InstallURL   string       `json:"installURL,omitempty"`
}
