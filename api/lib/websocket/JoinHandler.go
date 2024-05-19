package websocket

import (
	gameHelpers "api/lib/game"
	"errors"
)

func JoinHandler(message ClientMessage, client *Client, pool *ConnectionPool) error {
	// Add the player to the game.
	addError := gameHelpers.AddPlayerToGame(client.GameId, client.PlayerId)
	if addError != nil {
		return addError
	}

	// Broadcast the updated player list.
	broadcastError := BroadcastPlayersUpdate(client, pool)
	if broadcastError != nil {
		return errors.New("could not broadcast player list")
	}

	return nil
}
