package storages

type Song struct {
	ID    int    `json:"id"`
	Group string `json:"group"`
	Name  string `json:"song"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
