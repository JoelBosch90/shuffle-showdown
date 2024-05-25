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

	nextPlayerOrder := uint(0)

	if len(game.Rounds) > 0 {
		_, lastRound := LastRound(game.Rounds)

		lastPlayerOrder, findOrderError := findPlayerOrder(game.GamePlayers, lastRound.Player.Id)
		if findOrderError != nil {
			return uuid.UUID{}, findOrderError
		}

		nextPlayerOrder = (lastPlayerOrder + 1) % uint(len(game.GamePlayers))
	}

	nextPlayerId, findPlayerError := findPlayerIdByOrder(game.GamePlayers, nextPlayerOrder)
	if findPlayerError != nil {
		return uuid.UUID{}, findPlayerError
	}

	return nextPlayerId, nil
}
