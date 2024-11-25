package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TakuroBreath/song-library/internal/domain/models"
	"github.com/TakuroBreath/song-library/internal/storage"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func (s *SongService) GetSongVerses(group, song string, limit, offset int) ([]string, error) {
	s.log.Info("Getting song verses",
		slog.String("group", group),
		slog.String("song", song),
		slog.Int("limit", limit),
		slog.Int("offset", offset))

	verses, err := s.Storage.GetSongWithPagination(group, song, limit, offset)
	if err != nil {
		s.log.Error("Failed to get song verses",
			slog.String("group", group),
			slog.String("song", song),
			slog.Any("error", err))
		return nil, err
	}
	return verses, nil
}

func (s *SongService) GetSongs(filters map[string]interface{}, limit, offset int) ([]*models.Song, error) {
	s.log.Info("Getting filtered songs",
		slog.Any("filters", filters),
		slog.Int("limit", limit),
		slog.Int("offset", offset))

	songs, err := s.Storage.GetFilteredSongs(filters, limit, offset)
	if err != nil {
		s.log.Error("Failed to get filtered songs",
			slog.Any("filters", filters),
			slog.Any("error", err))
		return nil, err
	}
	return songs, nil
}

func (s *SongService) UpdateSong(id int, group, song, releaseDate, text, link *string) error {
	s.log.Info("Updating song",
		slog.Int("id", id),
		slog.Any("group", group),
		slog.Any("song", song))

	err := s.Storage.UpdateSong(id, group, song, releaseDate, text, link)
	if err != nil {
		s.log.Error("Failed to update song",
			slog.Int("id", id),
			slog.Any("error", err))
		return err
	}
	return nil
}

func (s *SongService) DeleteSong(group, song string) error {
	s.log.Info("Deleting song",
		slog.String("group", group),
		slog.String("song", song))

	err := s.Storage.DeleteSong(group, song)
	if err != nil {
		s.log.Error("Failed to delete song",
			slog.String("group", group),
			slog.String("song", song),
			slog.Any("error", err))
		return err
	}
	return nil
}

func (s *SongService) AddSongWithAPI(group, song string) (int, error) {
	s.log.Info("Adding song via API",
		slog.String("group", group),
		slog.String("song", song))

	reqUrl := fmt.Sprintf("%s/info?group=%s&song=%s", s.apiURL, url.QueryEscape(group), url.QueryEscape(song))
	resp, err := http.Get(reqUrl)
	if err != nil {
		s.log.Error("Failed to call external API",
			slog.String("url", reqUrl),
			slog.Any("error", err))
		return 0, fmt.Errorf("failed to call external API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.log.Error("External API returned non-OK status",
			slog.String("status", resp.Status))
		return 0, fmt.Errorf("API returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("Failed to read API response",
			slog.Any("error", err))
		return 0, fmt.Errorf("failed to read API response: %w", err)
	}

	var songDetail SongDetail
	if err := json.Unmarshal(body, &songDetail); err != nil {
		s.log.Error("Failed to parse API response",
			slog.Any("error", err))
		return 0, fmt.Errorf("failed to parse API response: %w", err)
	}

	songID, err := s.Storage.AddSong(group, song, songDetail.ReleaseDate, songDetail.Text, songDetail.Link)
	if err != nil {
		s.log.Error("Failed to save song in repository",
			slog.String("group", group),
			slog.String("song", song),
			slog.Any("error", err))

		if err == storage.ErrSongExists {
			return 0, storage.ErrSongExists
		}
		return 0, err
	}

	s.log.Info("Song added successfully",
		slog.Int("songID", songID),
		slog.String("group", group),
		slog.String("song", song))

	return songID, nil
}

func (s *SongService) GetID(group, song string) (int, error) {
	s.log.Info("Getting song ID",
		slog.String("group", group),
		slog.String("song", song))

	id, err := s.Storage.GetID(group, song)
	if err != nil {
		s.log.Error("Failed to get song ID",
			slog.String("group", group),
			slog.String("song", song),
			slog.Any("error", err))

		if errors.Is(err, storage.ErrSongNotFound) {
			return 0, storage.ErrSongNotFound
		}
		return 0, err
	}
	return id, nil
}
