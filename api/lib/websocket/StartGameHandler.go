package websocket

import (
	"api/database"
	"api/database/models"
	"api/lib/game"
	"errors"

	uuid "github.com/satori/go.uuid"
)

func runGame(gameId uuid.UUID) error {
	database := database.Get()

	var game models.Game
	loadGameError := database.Where("id = ?", gameId).First(&game).Error
	if loadGameError != nil {
		return errors.New("could not load game")
	}

	game.HasStarted = true
	saveGameError := database.Save(&game).Error
	if saveGameError != nil {
		return errors.New("could not save game")
	}

	return nil
}

func StartGameHandler(message ClientMessage, client *Client, pool *ConnectionPool) error {
	if !game.IsReadyToStart(client.GameId) {
		return errors.New("game is not ready to start")
	}

	randomizePlayerError := game.ShufflePlayers(client.GameId)
	if randomizePlayerError != nil {
		return errors.New("could not randomize player order")
	}

	awardInitialTracksError := game.AwardInitialTracks(client.GameId)
	if awardInitialTracksError != nil {
		return errors.New("could not award initial tracks")
	}

	createRoundError := game.CreateNextRound(client.GameId)
	if createRoundError != nil {
		return errors.New("could not create round")
	}

	setRunningError := runGame(client.GameId)
	if setRunningError != nil {
		return errors.New("could not start game")
	}

	broadcastError := BroadcastGameUpdate(client, pool)
	if broadcastError != nil {
		return errors.New("could not broadcast game update")
	}

	return nil
}
