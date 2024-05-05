package v1

import (
	"api/database"
	"api/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GameSettings struct {
	SongsToWin     uint `json:"songsToWin"`
	TitleRequired  bool `json:"titleRequired"`
	ArtistRequired bool `json:"artistRequired"`
}

type PatchGameInput struct {
	Settings GameSettings `json:"settings"`
	PlayerId string       `json:"playerId"`
}

func PatchGame(context *gin.Context) {
	var input PatchGameInput
	validationError := context.ShouldBindJSON(&input)
	if validationError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	gameId := context.Param("id")
	database := database.Get()
	var game models.Game

	cookie, cookieError := context.Cookie("playerSecret")
	if cookieError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cookie"})
		return
	}

	var player models.Player
	playerError := database.Where("secret = ? AND id = ?", cookie, input.PlayerId).First(&player).Error
	if playerError != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Game not owned by user"})
		return
	}

	gameError := database.Where("id = ?", gameId).First(&game).Error
	if gameError != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Game does not exist"})
		return
	}
	if game.Configured {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Game already configured"})
		return
	}

	database.Model(&game).Updates(models.Game{
		SongsToWin:     input.Settings.SongsToWin,
		TitleRequired:  input.Settings.TitleRequired,
		ArtistRequired: input.Settings.ArtistRequired,
		Configured:     true,
	})

	context.JSON(http.StatusOK, gin.H{"data": game})
}
