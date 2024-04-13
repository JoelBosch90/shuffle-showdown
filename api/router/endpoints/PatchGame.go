package endpoints

import (
	"api/database"
	"api/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PatchGameInput struct {
	SongsToWin     uint `json:"songsToWin"`
	TitleRequired  bool `json:"titleRequired"`
	ArtistRequired bool `json:"artistRequired"`
}

func PatchGame(context *gin.Context) {
	var input PatchGameInput
	validationError := context.ShouldBindJSON(&input)
	if validationError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	id := context.Param("id")
	database := database.Get()
	var game models.Game

	databaseError := database.Where("id = ?", id).First(&game).Error
	if databaseError != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if game.Configured {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Game already configured"})
		return
	}

	database.Model(&game).Updates(models.Game{
		SongsToWin:     input.SongsToWin,
		TitleRequired:  input.TitleRequired,
		ArtistRequired: input.ArtistRequired,
		Configured:     true,
	})

	context.JSON(http.StatusOK, gin.H{"data": game})
}
