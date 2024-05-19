package game

import (
	database "api/database"
	"api/database/models"
	"errors"

	uuid "github.com/satori/go.uuid"
)

func AddPlayerToGame(gameId uuid.UUID, playerId uuid.UUID) error {
	database := database.Get()
	var player models.Player
	var game models.Game

	playerError := database.Where("id = ?", playerId).First(&player).Error
	if playerError != nil {
		return errors.New("could not join game")
	}

	gameError := database.Preload("Players").Preload("BannedPlayers").Where("id = ?", gameId).First(&game).Error
	if gameError != nil {
		return errors.New("could not join game")
	}

	if IncludesPlayer(game.BannedPlayers, player) {
		return errors.New("player is banned")
	}

	// No need to add the player to the game if he's already added.
	if IncludesPlayer(game.Players, player) {
		return nil
	}

	// Add the player to the game.
	game.Players = append(game.Players, player)
	updateError := database.Save(&game).Error
	if updateError != nil {
		return errors.New("could not join game")
	}

	return nil
}
