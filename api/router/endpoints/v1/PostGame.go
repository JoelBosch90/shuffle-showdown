package v1

import (
	"api/database"
	gameHelpers "api/lib/game"
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

	// Extract the playlist ID from the input.
	playlistId := spotify.ExtractPlaylistId(input.Playlist)

	// Request the playlist information from Spotify.
	playlist, playlistError := spotify.RequestPlaylistInfo(playlistId, input.CountryCode)
	if playlistError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing playlist"})
		return
	}

	// Create a player and a game.
	database := database.Get()
	player, playerError := gameHelpers.CreatePlayer("", database)
	if playerError != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	}
	game, gameError := gameHelpers.CreateGame(playlist, player, database)
	if gameError != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	}

	// Set the player cookie.
	gameHelpers.SetPlayerCookie(context, player)

	context.JSON(http.StatusOK, gin.H{"game": game})
}
