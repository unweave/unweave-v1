package image

type Image struct {
	ID       string `json:"id"`
	Tag      string `json:"tag"`
	Repo     string `json:"repo"`
	Registry string `json:"registry"`
}
