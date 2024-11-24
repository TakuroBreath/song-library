package service

import (
	"encoding/json"
	"fmt"
	"github.com/TakuroBreath/song-library/internal/domain/models"
	"io"
	"net/http"
	"net/url"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func (s *SongService) GetSongVerses(group, song string, limit, offset int) ([]string, error) {
	return s.Storage.GetSongWithPagination(group, song, limit, offset)
}

func (s *SongService) GetSongs(filters map[string]interface{}, limit, offset int) ([]*models.Song, error) {
	return s.Storage.GetFilteredSongs(filters, limit, offset)
}

func (s *SongService) UpdateSong(id int, group, song, releaseDate, text, link *string) error {
	return s.Storage.UpdateSong(id, group, song, releaseDate, text, link)
}

func (s *SongService) DeleteSong(group, song string) error {
	return s.Storage.DeleteSong(group, song)
}

func (s *SongService) AddSongWithAPI(group, song string) (int, error) {
	reqUrl := fmt.Sprintf("%s/info?group=%s&song=%s", s.apiURL, url.QueryEscape(group), url.QueryEscape(song))
	resp, err := http.Get(reqUrl)
	if err != nil {
		return 0, fmt.Errorf("failed to call external API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read API response: %w", err)
	}

	var songDetail SongDetail
	if err := json.Unmarshal(body, &songDetail); err != nil {
		return 0, fmt.Errorf("failed to parse API response: %w", err)
	}

	songID, err := s.Storage.AddSong(group, song, songDetail.ReleaseDate, songDetail.Text, songDetail.Link)
	if err != nil {
		return 0, fmt.Errorf("failed to save song in repository: %w", err)
	}

	return songID, nil
}

func (s *SongService) GetID(group, song string) (int, error) {
	return s.Storage.GetID(group, song)
}
