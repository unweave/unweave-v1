package types

type Repository struct {
	Name string `json:"name,omitempty"`
	// URL must be cloneable by Git
	URL string `json:"url,omitempty"`
}

// GithubIntegrationGetResponse lists all repositories authenticated to Unweave,
// whether the GitHub application is installed,
// and the installation URL if applicable.
type GithubIntegrationGetResponse struct {
	Repositories   []Repository `json:"repositories,omitempty"`
	IsAppInstalled bool         `json:"isAppInstalled"`
	InstallURL     string       `json:"installURL,omitempty"`
}

type GithubIntegrationConnectRequest struct {
	Code        string `json:"code"`
	AccessToken string `json:"accessToken"`
}
