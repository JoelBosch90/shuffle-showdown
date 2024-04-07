package endpoints

import (
	"api/database"
	"api/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostGameInput struct {
	PlayListID uint `json:"playlist" binding:"required"`
}

func PostGame(context *gin.Context) {
	var input PostGameInput
	validationError := context.ShouldBindJSON(&input)
	if validationError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	database := database.Get()
	game := models.Game{PlayListID: input.PlayListID}
	databaseError := database.Create(&game).Error
	if databaseError != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	}

	context.JSON(http.StatusOK, gin.H{"data": game})
}
