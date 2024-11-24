package routes

import (
	"github.com/TakuroBreath/song-library/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupSongRoutes(router *gin.Engine, songHandler *handlers.SongHandler) {
	// Группа маршрутов для API песен
	songs := router.Group("/api/songs")
	{
		// GET /api/songs - получение списка песен с фильтрацией и пагинацией
		songs.GET("", songHandler.GetSongs)

		// GET /api/songs/verses - получение куплетов песни
		songs.GET("/verses", songHandler.GetSongVerses)

		// POST /api/songs - добавление новой песни
		songs.POST("", songHandler.AddSong)

		// PUT /api/songs - обновление информации о песне
		songs.PUT("", songHandler.UpdateSong)

		// DELETE /api/songs - удаление песни
		songs.DELETE("", songHandler.DeleteSong)
	}
}
