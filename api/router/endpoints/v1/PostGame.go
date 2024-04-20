package v1

import (
	"api/database"
	"api/lib/game"
	"api/lib/spotify"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostGameInput struct {
	PlayList    string `json:"playList" binding:"required"`
	CountryCode string `json:"countryCode" binding:"required"`
}

func PostGame(context *gin.Context) {
	var input PostGameInput
	validationError := context.ShouldBindJSON(&input)
	if validationError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	playListId := spotify.ExtractPlayListId(input.PlayList)
	info, playListError := spotify.RequestPlayListInfo(playListId, input.CountryCode)
	if playListError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing playlist"})
		return
	}

	database := database.Get()
	game, gameError := game.CreateGame(info, database)
	if gameError != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	}

	context.JSON(http.StatusOK, gin.H{"data": game})
}