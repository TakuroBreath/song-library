package handlers

import (
	"github.com/TakuroBreath/song-library/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type SongHandler struct {
	songService *service.SongService
}

func NewSongHandler(songService *service.SongService) *SongHandler {
	return &SongHandler{songService: songService}
}

func (h *SongHandler) GetSongs(c *gin.Context) {
	filters := map[string]interface{}{}

	if group := c.Query("group"); group != "" {
		filters["group"] = group
	}
	if song := c.Query("song"); song != "" {
		filters["song"] = song
	}
	if releaseDate := c.Query("release_date"); releaseDate != "" {
		filters["release_date"] = releaseDate
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	songs, err := h.songService.GetSongs(filters, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, songs)
}

func (h *SongHandler) GetSongVerses(c *gin.Context) {
	group := c.Query("group")
	song := c.Query("song")

	if group == "" || song == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group and song are required"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	verses, err := h.songService.GetSongVerses(group, song, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, verses)
}
