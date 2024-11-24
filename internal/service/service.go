package service

import "github.com/TakuroBreath/song-library/internal/storage/postgresql"

type SongService struct {
	Storage *postgresql.Storage
	apiURL  string
}

func NewSongService(storage *postgresql.Storage, apiURL string) *SongService {
	return &SongService{Storage: storage, apiURL: apiURL}
}
