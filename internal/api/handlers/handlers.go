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

type SongAddRequest struct {
	Group string `json:"group" binding:"required,min=1,max=255"`
	Song  string `json:"song"  binding:"required,min=1,max=255"`
}

type SongUpdateRequest struct {
	Group       *string `json:"group,omitempty" binding:"omitempty,min=1,max=255"`
	Song        *string `json:"song,omitempty"  binding:"omitempty,min=1,max=255"`
	ReleaseDate *string `json:"release_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
	Text        *string `json:"text,omitempty"`
	Link        *string `json:"link,omitempty" binding:"omitempty,url"`
}

// GetSongs godoc
// @Summary      Get songs list
// @Description  Get songs with filtering and pagination
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        group query string false "Filter by group name"
// @Param        song query string false "Filter by song name"
// @Param        release_date query string false "Filter by release date (YYYY-MM-DD)"
// @Param        limit query int false "Limit number of records" default(10)
// @Param        offset query int false "Offset for pagination" default(0)
// @Success      200  {array}   models.Song
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /songs [get]
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

// GetSongVerses godoc
// @Summary      Get song verses
// @Description  Get verses of a specific song with pagination
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        group query string true "Group name"
// @Param        song query string true "Song name"
// @Param        limit query int false "Limit number of verses" default(5)
// @Param        offset query int false "Offset for pagination" default(0)
// @Success      200  {array}   string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /songs/verses [get]
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

// DeleteSong godoc
// @Summary      Delete song
// @Description  Delete existing song
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        group query string true "Group name"
// @Param        song query string true "Song name"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /songs [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	group := c.Query("group")
	song := c.Query("song")

	if group == "" || song == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group and song are required"})
		return
	}

	err := h.songService.DeleteSong(group, song)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "song deleted successfully"})
}

// AddSong godoc
// @Summary      Add new song
// @Description  Add a new song with details from external API
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        request body SongAddRequest true "Song details"
// @Success      201  {object}  map[string]int
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /songs [post]
func (h *SongHandler) AddSong(c *gin.Context) {
	var request SongAddRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	songID, err := h.songService.AddSongWithAPI(request.Group, request.Song)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": songID})
}

// UpdateSong godoc
// @Summary      Update song
// @Description  Update existing song details
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        group query string true "Group name"
// @Param        song query string true "Song name"
// @Param        request body SongUpdateRequest true "Song update details"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /songs [put]
func (h *SongHandler) UpdateSong(c *gin.Context) {
	group := c.Query("group")
	song := c.Query("song")

	if group == "" || song == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group and song are required"})
		return
	}
	var request SongUpdateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.songService.GetID(group, song)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.songService.UpdateSong(id, request.Group, request.Song, request.ReleaseDate, request.Text, request.Link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "message": "song updated successfully"})
}
