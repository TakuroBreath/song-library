package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/TakuroBreath/song-library/internal/domain/models"
	_ "github.com/lib/pq"
	"strings"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(psqlInfo string) (*Storage, error) {
	const op = "storage.postgresql.NewStorage"

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) AddSong(group, song, releaseDate, text, link string) (int, error) {
	const op = "storage.postgresql.AddSong"

	var id int

	err := s.db.QueryRow(`
    						INSERT INTO songs ("group", song, release_date, text, link) 
    						VALUES ($1, $2, $3, $4, $5)
    						RETURNING id
							`, group, song, releaseDate, text, link,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) UpdateSong(id int, group, song, releaseDate, text, link string) error {
	const op = "storage.postgresql.UpdateSong"

	_, err := s.db.Exec(`
       						UPDATE songs 
       						SET "group" = $1, song = $2, release_date = $3, text = $4, link = $5
       						WHERE id = $6
   							`, group, song, releaseDate, text, link, id)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteSong(group, song string) error {
	const op = "storage.postgresql.DeleteSong"

	result, err := s.db.Exec(`
       DELETE FROM songs 
       WHERE "group" = $1 AND song = $2
   `, group, song)

	if err != nil {
		return fmt.Errorf("%s: execute delete: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: song '%s' by '%s' not found", op, song, group)
	}

	return nil
}

func (s *Storage) GetSongWithPagination(group, song string, limit, offset int) ([]string, error) {
	const op = "storage.postgresql.GetSongWithPagination"

	var text string

	err := s.db.QueryRow(`
		SELECT text 
		FROM songs 
		WHERE "group" = $1 AND song = $2
	`, group, song).Scan(&text)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	verses := strings.Split(text, "\n\n")

	if offset >= len(verses) {
		return nil, nil
	}

	end := offset + limit
	if end > len(verses) {
		end = len(verses)
	}

	return verses[offset:end], nil
}

func (s *Storage) GetFilteredSongs(filters map[string]interface{}, limit, offset int) ([]*models.Song, error) {
	const op = "storage.postgresql.GetFilteredSongs"

	query := `SELECT id, "group", song, release_date, text, link FROM songs`
	var conditions []string
	var args []interface{}
	argIndex := 1

	for field, value := range filters {
		conditions = append(conditions, fmt.Sprintf(`%s = $%d`, field, argIndex))
		args = append(args, value)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var songs []*models.Song

	for rows.Next() {
		var songDetail models.Song
		err := rows.Scan(&songDetail.ID, &songDetail.Group, &songDetail.Song, &songDetail.ReleaseDate, &songDetail.Text, &songDetail.Link)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		songs = append(songs, &songDetail)
	}

	return songs, nil
}
