package websocket

import (
	"api/database"
	"api/database/models"
	gameHelpers "api/lib/game"
	"errors"
)

func StartGameHandler(message ClientMessage, client *Client, pool *ConnectionPool) error {
	database := database.Get()

	randomizePlayerError := gameHelpers.ShufflePlayers(client.GameId)
	if randomizePlayerError != nil {
		return errors.New("could not randomize player order")
	}

	setRunningError := database.Save(&models.Game{Id: client.GameId, IsRunning: true}).Error
	if setRunningError != nil {
		return errors.New("could not start game")
	}

	broadcastError := BroadcastGameUpdate(client, pool)
	if broadcastError != nil {
		return errors.New("could not broadcast game update")
	}

	return nil
}
