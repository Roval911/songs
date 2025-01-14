package postgres

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"songs/internal/storages"
)

func (s *PostgresStorage) AddLyrics(songID int, line string) error {
	query := `INSERT INTO song_lyrics (song_id, lyrics_line) VALUES ($1, $2)`
	_, err := s.db.Exec(query, songID, line)
	if err != nil {
		s.logger.Printf("Ошибка при добавлении строки текста песни (songID: %d): %v", songID, err)
		return err
	}
	return nil
}

func (s *PostgresStorage) UpdateSongPartial(id int, updates map[string]interface{}) error {
	query := `UPDATE songs SET `
	params := []interface{}{}
	i := 1

	for key, value := range updates {
		if i > 1 {
			query += ", "
		}
		query += fmt.Sprintf(`"%s" = $%d`, key, i)
		params = append(params, value)
		i++
	}
	query += " WHERE id = $1"
	params = append(params, id)

	_, err := s.db.Exec(query, params...)
	if err != nil {
		s.logger.Printf("Ошибка при частичном обновлении песни (ID: %d): %v", id, err)
		return err
	}
	return nil
}

func (s *PostgresStorage) GetSongs(group string, song string, page int, limit int) ([]storages.Song, error) {
	offset := (page - 1) * limit
	query := `
        SELECT s.id, g.name AS group_name, s.name, s.release_date, s.link
        FROM songs s
        JOIN groups g ON s.group_id = g.id
        WHERE g.name ILIKE $1 AND s.name ILIKE $2
        LIMIT $3 OFFSET $4;
    `
	rows, err := s.db.Query(query, "%"+group+"%", "%"+song+"%", limit, offset)
	if err != nil {
		s.logger.Printf("Ошибка при получении песен: %v", err)
		return nil, err
	}
	defer rows.Close()

	var songs []storages.Song
	for rows.Next() {
		var song storages.Song
		if err := rows.Scan(&song.ID, &song.Group, &song.Name, &song.ReleaseDate, &song.Link); err != nil {
			s.logger.Printf("Ошибка при сканировании результата: %v", err)
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}

func (s *PostgresStorage) GetLyrics(songID int, page int, limit int) ([]string, error) {
	offset := (page - 1) * limit
	query := `
        SELECT lyrics_line
        FROM song_lyrics
        WHERE song_id = $1
        LIMIT $2 OFFSET $3;
    `
	rows, err := s.db.Query(query, songID, limit, offset)
	if err != nil {
		s.logger.Printf("Ошибка при получении текста песни: %v", err)
		return nil, err
	}
	defer rows.Close()

	var lyrics []string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			s.logger.Printf("Ошибка при сканировании текста: %v", err)
			return nil, err
		}
		lyrics = append(lyrics, line)
	}
	return lyrics, nil
}

func (s *PostgresStorage) DeleteSong(id int) error {
	query := `DELETE FROM songs WHERE id = $1`
	_, err := s.db.Exec(query, id)
	if err != nil {
		s.logger.Printf("Ошибка при удалении песни (ID: %d): %v", id, err)
		return err
	}
	s.logger.Printf("Успешно удалена песня (ID: %d)", id)
	return nil
}

func (s *PostgresStorage) UpdateSong(id int, song storages.Song) error {
	query := `UPDATE songs SET "group" = $1, song = $2 WHERE id = $3`
	_, err := s.db.Exec(query, song.Group, song.Name, id)
	if err != nil {
		s.logger.Printf("Ошибка при обновлении песни (ID: %d): %v", id, err)
		return err
	}
	s.logger.Printf("Успешно обновлена песня (ID: %d)", id)
	return nil
}

func (s *PostgresStorage) AddSong(song storages.Song) error {
	query := `
        INSERT INTO songs (group_id, name, release_date, link)
        VALUES ((SELECT id FROM groups WHERE name = $1), $2, $3, $4);
    `
	_, err := s.db.Exec(query, song.Group, song.Name, song.ReleaseDate, song.Link)
	if err != nil {
		s.logger.Printf("Ошибка при добавлении песни (группа: %s, песня: %s): %v", song.Group, song.Name, err)
		return err
	}
	s.logger.Printf("Песня успешно добавлена (группа: %s, песня: %s)", song.Group, song.Name)
	return nil
}

func (s *PostgresStorage) CreateIndexes() error {
	if s.logger == nil {
		s.logger = logrus.New() // Если логгер не был инициализирован, создаем новый
	}
	_, err := s.db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_group_name ON songs(group_id);
		CREATE INDEX IF NOT EXISTS idx_song_name ON songs(name);
		CREATE INDEX IF NOT EXISTS idx_song_id ON song_lyrics(song_id);
	`)
	if err != nil {
		s.logger.Printf("Ошибка при создании индексов: %v", err)
		return err
	}
	s.logger.Println("Индексы успешно созданы.")
	return nil
}
