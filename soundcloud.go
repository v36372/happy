package happy

type SCResolveResponse struct {
	Kind       string  `json:"kind"`
	Title      string  `json:"title"`
	URI        string  `json:"uri"`
	ArtworkURL string  `json:"artwork_url"`
	Tracks     []Track `json:"tracks"`
}

type Track struct {
	Title      string `json:"title"`
	URI        string `json:"uri"`
	ArtworkURL string `json:"artwork_url"`
}
