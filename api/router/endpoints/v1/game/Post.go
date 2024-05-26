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

	playlistId := spotify.ExtractPlaylistId(input.Playlist)
	playlist, playlistError := spotify.RequestPlaylistInfo(playlistId, input.CountryCode)
	if playlistError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing playlist"})
		return
	}

	database := database.Get()

	player := &models.Player{}
	if input.PlayerId != uuid.Nil {
		getPlayerError := database.Where("id = ?", input.PlayerId).First(player).Error
		if getPlayerError != nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		}
	} else {
		var createPlayerError error
		player, createPlayerError = gameHelpers.CreatePlayer("", database)
		if createPlayerError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
	}

	game, gameError := gameHelpers.CreateGame(playlist, *player, database)
	if gameError != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	}

	gameHelpers.SetPlayerCookie(context, *player)
	context.JSON(http.StatusOK, gin.H{"game": game})
}
