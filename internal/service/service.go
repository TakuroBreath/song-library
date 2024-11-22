package service

import "github.com/TakuroBreath/song-library/internal/storage/postgresql"

type SongService struct {
	storage *postgresql.Storage
	apiURL  string
}

func NewSongService(storage *postgresql.Storage, apiURL string) *SongService {
	return &SongService{storage: storage, apiURL: apiURL}
}
