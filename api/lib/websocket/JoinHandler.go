package websocket

import (
	gameHelpers "api/lib/game"
	"errors"

	uuid "github.com/satori/go.uuid"
)

func JoinHandler(message ClientMessage, client *Client, pool *ConnectionPool) error {
	// Check if the player has already identified himself.
	playerId, uuidError := uuid.FromString(message.PlayerId.String())
	if uuidError != nil || client.PlayerId != playerId {
		return errors.New("player identification failed")
	}

	// Add the player to the game.
	addError := gameHelpers.AddPlayerToGame(client.GameId, playerId)
	if addError != nil {
		return errors.New("could not join game")
	}

	// Broadcast the updated player list.
	broadcastError := BroadcastPlayersUpdate(client, pool)
	if broadcastError != nil {
		return errors.New("could not broadcast player list")
	}

	return nil
}
