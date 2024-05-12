package websocket

import (
	gameHelpers "api/lib/game"
	"encoding/json"
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

	// Get the player names to send to the client.
	names, namesError := gameHelpers.GetPlayerNames(client.GameId)
	if namesError != nil {
		return errors.New("could not read player names")
	}

	// Convert to JSON.
	namesJson, jsonError := json.Marshal(&PlayerNames{Players: names})
	if jsonError != nil {
		return errors.New("could not read player names")
	}

	// Broadcast the joined message to all plaers in the game.
	pool.Broadcast <- ServerMessage{
		Type:    ServerMessageTypeJoined,
		Content: string(namesJson),
		GameId:  client.GameId,
	}

	return nil
}
