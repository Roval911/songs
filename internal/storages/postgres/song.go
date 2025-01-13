package postgres

import (
	_ "github.com/lib/pq"
	"songs/internal/storages"
)

func (s *PostgresStorage) GetSongs(group string, song string, page int, limit int) ([]storages.Song, error) {
	offset := (page - 1) * limit
	query := `SELECT id, "group", song FROM songs WHERE "group" ILIKE $1 AND song ILIKE $2 LIMIT $3 OFFSET $4`
	rows, err := s.db.Query(query, "%"+group+"%", "%"+song+"%", limit, offset)
	if err != nil {
		s.logger.Printf("Ошибка при получении песен (группа: %s, песня: %s, страница: %d, лимит: %d): %v", group, song, page, limit, err)
		return nil, err
	}
	defer rows.Close()

	var songs []storages.Song
	for rows.Next() {
		var song storages.Song
		if err := rows.Scan(&song.ID, &song.Group, &song.Name); err != nil {
			// Логирование ошибки при чтении строки результата
			s.logger.Printf("Ошибка при сканировании строки результата: %v", err)
			return nil, err
		}
		songs = append(songs, song)
	}
	s.logger.Printf("Успешно получено %d песен (группа: %s, песня: %s, страница: %d, лимит: %d)", len(songs), group, song, page, limit)
	return songs, nil
}

func (s *PostgresStorage) GetLyrics(songID int, page int, limit int) ([]string, error) {
	offset := (page - 1) * limit
	query := `SELECT lyrics FROM song_lyrics WHERE song_id = $1 LIMIT $2 OFFSET $3`
	rows, err := s.db.Query(query, songID, limit, offset)
	if err != nil {
		s.logger.Printf("Ошибка при получении текста песни (songID: %d, страница: %d, лимит: %d): %v", songID, page, limit, err)
		return nil, err
	}
	defer rows.Close()

	var lyrics []string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			s.logger.Printf("Ошибка при сканировании строки текста песни (songID: %d): %v", songID, err)
			return nil, err
		}
		lyrics = append(lyrics, line)
	}
	s.logger.Printf("Успешно получено %d строк текста песни (songID: %d, страница: %d, лимит: %d)", len(lyrics), songID, page, limit)
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
	query := `INSERT INTO songs ("group", song) VALUES ($1, $2)`
	_, err := s.db.Exec(query, song.Group, song.Name)
	if err != nil {
		s.logger.Printf("Ошибка при добавлении новой песни (группа: %s, песня: %s): %v", song.Group, song.Name, err)
		return err
	}
	s.logger.Printf("Успешно добавлена новая песня (группа: %s, песня: %s)", song.Group, song.Name)
	return nil
}
