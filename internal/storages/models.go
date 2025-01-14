package storages

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Song struct {
	ID          int    `json:"id"`
	GroupID     int    `json:"group_id"`
	Group       string `json:"group"`
	Name        string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Link        string `json:"link"`
}

type SongDetail struct {
	ReleaseDate string   `json:"releaseDate"`
	Text        []string `json:"text"`
	Link        string   `json:"link"`
}
