package game

import (
	database "api/database"
	"api/database/models"
	"errors"

	uuid "github.com/satori/go.uuid"
)

func findPlayerIdByOrder(gamePlayers []models.GamePlayer, order uint) (uuid.UUID, error) {
	for _, gamePlayer := range gamePlayers {
		if gamePlayer.Order == order {
			return gamePlayer.PlayerId, nil
		}
	}
	return uuid.UUID{}, errors.New("player not found")
}

func findPlayerOrder(gamePlayers []models.GamePlayer, playerId uuid.UUID) (uint, error) {
	for _, gamePlayer := range gamePlayers {
		if gamePlayer.PlayerId == playerId {
			return gamePlayer.Order, nil
		}
	}
	return 0, errors.New("player not found")
}

func FindNextPlayerId(gameId uuid.UUID) (uuid.UUID, error) {
	database := database.Get()
	var game models.Game

	loadGameError := database.Preload("GamePlayers").Preload("Players").Preload("Rounds.Player").Where("id = ?", gameId).First(&game).Error
	if loadGameError != nil {
		return uuid.UUID{}, errors.New("could not load game")
	}
	if len(game.GamePlayers) < 2 {
		return uuid.UUID{}, errors.New("too few players")
	}

	nextPlayerId := game.GamePlayers[0].PlayerId

	if len(game.Rounds) > 0 {
		lastRound := game.Rounds[len(game.Rounds)-1]
		lastPlayerId := lastRound.Player.Id

		lastPlayerOrder, findOrderError := findPlayerOrder(game.GamePlayers, lastPlayerId)
		if findOrderError != nil {
			return uuid.UUID{}, findOrderError
		}

		var findPlayerError error
		nextPlayerOrder := (lastPlayerOrder + 1) % uint(len(game.GamePlayers))
		nextPlayerId, findPlayerError = findPlayerIdByOrder(game.GamePlayers, nextPlayerOrder)
		if findPlayerError != nil {
			return uuid.UUID{}, findPlayerError
		}
	}

	return nextPlayerId, nil
}
