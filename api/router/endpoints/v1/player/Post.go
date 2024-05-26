package player

import (
	"api/database"
	gameHelpers "api/lib/game"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostPlayerInput struct {
	Name string `json:"name"`
}

func Post(context *gin.Context) {
	var input PostPlayerInput
	validationError := context.ShouldBindJSON(&input)
	if validationError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Create a player.
	database := database.Get()
	player, playerError := gameHelpers.CreatePlayer(input.Name, database)
	if playerError != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	}

	// Set the player cookie.
	gameHelpers.SetPlayerCookie(context, *player)

	context.JSON(http.StatusOK, gin.H{"player": *player})
}
