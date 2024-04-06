package endpoints

import (
	"api/database"
	"api/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PatchGameInput struct {
	PlayListID uint `json:"playlist"`
}

func PatchGame(context *gin.Context) {
	var input PatchGameInput
	validationError := context.ShouldBindJSON(&input)
	if validationError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
		return
	}

	id := context.Param("id")
	database := database.Get()
	var game models.Game

	databaseError := database.Where("id = ?", id).First(&game).Error
	if databaseError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": databaseError.Error()})
		return
	}

	database.Model(&game).Updates(models.Game{PlayListID: input.PlayListID})

	context.JSON(http.StatusOK, gin.H{"data": game})
}
