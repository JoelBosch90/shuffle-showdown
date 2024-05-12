package game

import (
	database "api/database"
	"api/database/models"

	uuid "github.com/satori/go.uuid"
)

func containsPlayer(haystack []models.Player, needle models.Player) bool {
	for _, player := range haystack {
		if player.Id == needle.Id {
			return true
		}
	}
	return false
}

func AddPlayerToGame(gameId uuid.UUID, playerId uuid.UUID) error {
	database := database.Get()
	var player models.Player
	var game models.Game

	playerError := database.Where("id = ?", playerId).First(&player).Error
	if playerError != nil {
		return playerError
	}

	gameError := database.Preload("Players").Where("id = ?", gameId).First(&game).Error
	if gameError != nil {
		return gameError
	}

	// No need to add the player to the game if he's already added.
	if containsPlayer(game.Players, player) {
		return nil
	}

	// Add the player to the game.
	game.Players = append(game.Players, player)
	updateError := database.Save(&game).Error

	return updateError
}
