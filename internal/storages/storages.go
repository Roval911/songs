package storages

type Storages interface {
	GetSongs(group string, song string, page int, limit int) ([]Song, error)
	GetLyrics(songID int, page int, limit int) ([]string, error)
	DeleteSong(id int) error
	UpdateSong(id int, song Song) error
	AddSong(song Song) error
}
