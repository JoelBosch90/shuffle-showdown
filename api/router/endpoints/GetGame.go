package endpoints

import (
	"api/database"
	"api/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGame(context *gin.Context) {
	id := context.Param("id")
	database := database.Get()
	var game models.Game

	databaseError := database.Where("id = ?", id).First(&game).Error
	if databaseError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": databaseError.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": game})
}