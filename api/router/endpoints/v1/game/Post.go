package game

import (
	"api/database"
	"api/database/models"
	gameHelpers "api/lib/game"
	"api/lib/spotify"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type PostGameInput struct {
	Playlist    string    `json:"playlist" binding:"required"`
	CountryCode string    `json:"countryCode" binding:"required"`
	PlayerId    uuid.UUID `json:"playerId"`
}

func Post(context *gin.Context) {
	var input PostGameInput
	validationError := context.ShouldBindJSON(&input)
	if validationError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	database := database.Get()
	playlistId := spotify.ExtractPlaylistId(input.Playlist)

	playlist, playlistError := spotify.GetRecentPlaylist(playlistId, input.CountryCode)
	if playlistError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error getting playlist"})
		return
	}

	player := models.Player{}
	if input.PlayerId != uuid.Nil {
		database.Where("id = ?", input.PlayerId).First(&player)
	}

	if player.Id == uuid.Nil {
		playerPointer, createPlayerError := gameHelpers.CreatePlayer("", database)
		player = *playerPointer
		if createPlayerError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
	}

	game, gameError := gameHelpers.CreateGame(playlist, player, database)
	if gameError != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	}

	gameHelpers.SetPlayerCookie(context, player)
	context.JSON(http.StatusOK, gin.H{"game": game})
}
