package v1

import (
	"api/database"
	"api/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PatchPlayerNameInput struct {
	PlayerName string `json:"playerName"`
}

func PatchPlayerName(context *gin.Context) {
	var input PatchPlayerNameInput
	validationError := context.ShouldBindJSON(&input)
	if validationError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	cookie, cookieError := context.Cookie("playerSecret")
	if cookieError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cookie"})
		return
	}

	var player models.Player
	database := database.Get()
	playerId := context.Param("id")
	playerError := database.Where("secret = ? AND id = ?", cookie, playerId).First(&player).Error
	if playerError != nil {
		context.JSON(http.StatusForbidden, gin.H{"error": "Player not authenticated"})
		return
	}

	// Add the player name.
	database.Model(&player).Updates(models.Player{
		Name: input.PlayerName,
	})

	context.JSON(http.StatusOK, gin.H{"data": player})
}
