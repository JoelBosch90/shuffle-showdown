package websocket

import (
	"api/database"
	"api/database/models"
	"errors"
)

func StartGameHandler(message ClientMessage, client *Client, pool *ConnectionPool) error {
	database := database.Get()

	updateError := database.Model(&models.Game{}).Where("id = ?", client.GameId).Update("running", true).Error
	if updateError != nil {
		return errors.New("could not start game")
	}

	broadcastError := BroadcastGameUpdate(client, pool)
	if broadcastError != nil {
		return errors.New("could not broadcast game update")
	}

	return nil
}
