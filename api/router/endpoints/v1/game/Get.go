package game

import (
	"api/database"
	"api/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get(context *gin.Context) {
	id := context.Param("id")
	database := database.Get()
	var game models.Game

	databaseError := database.Preload("Owner").Preload("Playlist").Preload("Players").Preload("Rounds").Where("id = ?", id).First(&game).Error
	if databaseError != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"game": game})
}
