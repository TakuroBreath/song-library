package service

import (
	"github.com/TakuroBreath/song-library/internal/storage/postgresql"
	"log/slog"
)

type SongService struct {
	Storage *postgresql.Storage
	apiURL  string
	log     *slog.Logger
}

func NewSongService(storage *postgresql.Storage, apiURL string, log *slog.Logger) *SongService {
	return &SongService{Storage: storage, apiURL: apiURL, log: log}
}
