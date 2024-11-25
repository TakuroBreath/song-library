package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/TakuroBreath/song-library/internal/domain/models"
	"github.com/TakuroBreath/song-library/internal/storage"
	_ "github.com/lib/pq"
	"log/slog"
	"strings"
)

type Storage struct {
	db  *sql.DB
	log *slog.Logger
}

func NewStorage(psqlInfo string, log *slog.Logger) (*Storage, error) {
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
		db:  db,
		log: log,
	}, nil
}

func (s *Storage) AddSong(group, song, releaseDate, text, link string) (int, error) {
	const op = "storage.postgresql.AddSong"

	var exists bool
	err := s.db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM songs WHERE "group" = $1 AND song = $2
        )
    `, group, song).Scan(&exists)

	if err != nil {
		return 0, fmt.Errorf("%s: check song existence: %w", op, err)
	}

	if exists {
		s.log.Warn("Attempt to add existing song",
			slog.String("group", group),
			slog.String("song", song))
		return 0, storage.ErrSongExists
	}

	var id int
	err = s.db.QueryRow(`
        INSERT INTO songs ("group", song, release_date, text, link) 
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, group, song, releaseDate, text, link).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("Song added successfully",
		slog.Int("id", id),
		slog.String("group", group),
		slog.String("song", song))

	return id, nil
}

func (s *Storage) UpdateSong(id int, group, song *string, releaseDate *string, text *string, link *string) error {
	const op = "storage.postgresql.UpdateSong"

	_, err := s.db.Exec(`
        UPDATE songs 
        SET "group" = COALESCE($1, "group"), 
            song = COALESCE($2, song),
            release_date = COALESCE($3, release_date), 
            text = COALESCE($4, text),
            link = COALESCE($5, link)
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

	// Заменяем экранированные переводы строк на реальные
	text = strings.ReplaceAll(text, "\\n", "\n")

	// Разделяем текст на куплеты
	verses := strings.Split(text, "\n\n")

	// Если получился только один куплет, пробуем разбить по одинарным переводам строк
	if len(verses) <= 1 {
		verses = strings.Split(text, "\n")
	}

	// Очищаем пустые строки и форматируем текст
	var cleanVerses []string
	for _, verse := range verses {
		verse = strings.TrimSpace(verse)
		if verse != "" {
			// Заменяем переводы строк на пробелы для форматированного вывода
			formattedVerse := strings.ReplaceAll(verse, "\n", " ")
			cleanVerses = append(cleanVerses, formattedVerse)
		}
	}

	// Проверяем границы пагинации
	if offset >= len(cleanVerses) {
		return []string{}, nil
	}

	end := offset + limit
	if end > len(cleanVerses) {
		end = len(cleanVerses)
	}

	return cleanVerses[offset:end], nil
}

func (s *Storage) GetFilteredSongs(filters map[string]interface{}, limit, offset int) ([]*models.Song, error) {
	const op = "storage.postgresql.GetFilteredSongs"

	query := `SELECT id, "group", song, release_date, text, link FROM songs`
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Карта для правильного экранирования имен полей
	fieldNames := map[string]string{
		"group":        `"group"`,
		"song":         "song",
		"release_date": "release_date",
		"text":         "text",
		"link":         "link",
	}

	for field, value := range filters {
		if quotedField, ok := fieldNames[field]; ok {
			conditions = append(conditions, fmt.Sprintf(`%s = $%d`, quotedField, argIndex))
			args = append(args, value)
			argIndex++
		}
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
func (s *Storage) GetID(group, song string) (int, error) {
	const op = "storage.postgresql.GetID"

	var id int

	err := s.db.QueryRow(`
        SELECT id 
        FROM songs 
        WHERE "group" = $1 AND song = $2
    `, group, song).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		s.log.Warn("Song not found",
			slog.String("group", group),
			slog.String("song", song))
		return 0, storage.ErrSongNotFound
	}

	if err != nil {
		s.log.Error("Failed to get song ID",
			slog.String("group", group),
			slog.String("song", song),
			slog.Any("error", err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
