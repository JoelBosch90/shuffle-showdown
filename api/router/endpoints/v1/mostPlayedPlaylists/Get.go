package game

import (
	"api/database"
	"api/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PlaylistWithGameCount struct {
	Playlist    models.Playlist `json:"playlist" gorm:"embedded"`
	GamesPlayed int             `json:"gamesPlayed"`
}

func Get(context *gin.Context) {
	database := database.Get()
	var playlists []PlaylistWithGameCount

	databaseError := database.Table("playlists").
		Select("playlists.*, COUNT(games.id) AS games_played").
		Joins("LEFT JOIN games ON games.playlist_id = playlists.id").
		Group("playlists.id").
		Order("games_played DESC").
		Limit(5).
		Find(&playlists).Error
	if databaseError != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Error finding playlists"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"playlists": playlists})
}
