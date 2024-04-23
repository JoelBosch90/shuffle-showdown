package v1

import (
	"api/database"
	"api/lib/game"
	"api/lib/spotify"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostGameInput struct {
	Playlist    string `json:"playlist" binding:"required"`
	CountryCode string `json:"countryCode" binding:"required"`
}

func PostGame(context *gin.Context) {
	var input PostGameInput
	validationError := context.ShouldBindJSON(&input)
	if validationError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	playlistId := spotify.ExtractPlaylistId(input.Playlist)
	playlist, playlistError := spotify.RequestPlaylistInfo(playlistId, input.CountryCode)
	if playlistError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing playlist"})
		return
	}

	database := database.Get()
	game, gameError := game.CreateGame(playlist, database)
	if gameError != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	}

	context.JSON(http.StatusOK, gin.H{"data": game})
}
