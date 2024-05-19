package game

import (
	database "api/database"
	"api/database/models"
	"errors"

	uuid "github.com/satori/go.uuid"
)

func BanPlayerFromGame(gameId uuid.UUID, playerId uuid.UUID) error {
	database := database.Get()
	var player models.Player
	var game models.Game

	playerError := database.Where("id = ?", playerId).First(&player).Error
	if playerError != nil {
		return playerError
	}

	gameError := database.Preload("BannedPlayers").Preload("Players").Where("id = ?", gameId).First(&game).Error
	if gameError != nil {
		return gameError
	}

	if !IncludesPlayer(game.BannedPlayers, player) {
		game.BannedPlayers = append(game.BannedPlayers, player)
		updateError := database.Save(&game).Error
		if updateError != nil {
			return errors.New("could not ban player")
		}
	}

	if IncludesPlayer(game.Players, player) {
		association := database.Model(&game).Association("Players")
		associationError := association.Delete(&player).Error
		if associationError != nil {
			return errors.New("could not kick player")
		}
	}

	return nil
}
