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
	ResourceID ResourceID `json:"resourceId"`
}

type ResourceID struct {
	VideoID string `json:"videoId"`
}

type Thumbnails map[string]*Thumbnail

type Thumbnail struct {
	URL string `json:"url"`
}
