package happy

type YTVideoListResponse struct {
	Items []Item `json:"items"`
}

type Item struct {
	ID      string  `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type Snippet struct {
	Title      string     `json:"title"`
	Thumbnails Thumbnails `json:"thumbnails"`
}

type Thumbnails struct {
	Default Thumbnail `json:"default"`
}

type Thumbnail struct {
	URL string `json:"url"`
}
