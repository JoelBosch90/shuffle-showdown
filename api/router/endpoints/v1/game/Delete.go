package game

import (
	"api/database"
	"api/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Delete(context *gin.Context) {
	id := context.Param("id")
	database := database.Get()
	var game models.Game

	databaseError := database.Where("id = ?", id).First(&game).Error
	if databaseError != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	database.Delete(&game)

	context.JSON(http.StatusOK, gin.H{"deleted": true})
}
