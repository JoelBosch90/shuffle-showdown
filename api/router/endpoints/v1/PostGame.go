package v1

import (
	"api/database"
	"api/database/models"
	"api/lib"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostGameInput struct {
	PlayListId string `json:"playListId" binding:"required"`
}

func PostGame(context *gin.Context) {
	var input PostGameInput
	validationError := context.ShouldBindJSON(&input)
	if validationError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	log.Println("Creating game with playlist link: ", input.PlayListId)
	lib.RequestSpotifyPlayListInfo(input.PlayListId)

	database := database.Get()
	game := models.Game{PlayListId: input.PlayListId}
	databaseError := database.Create(&game).Error
	if databaseError != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	}

	context.JSON(http.StatusOK, gin.H{"data": game})
}
