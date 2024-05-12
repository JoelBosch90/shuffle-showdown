package game

import (
	database "api/database"
	"api/database/models"

	uuid "github.com/satori/go.uuid"
)

func RemovePlayerFromGame(gameId uuid.UUID, playerId uuid.UUID) error {
	database := database.Get()
	var playerToRemove models.Player
	var game models.Game

	playerError := database.Where("id = ?", playerId).First(&playerToRemove).Error
	if playerError != nil {
		return playerError
	}

	gameError := database.Where("id = ?", gameId).First(&game).Error
	if gameError != nil {
		return gameError
	}

	// Remove the player from the game.
	association := database.Model(&game).Association("Players")
	associationError := association.Delete(&playerToRemove).Error

	return associationError
}
