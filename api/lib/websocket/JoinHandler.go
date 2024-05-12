package websocket

import (
	"api/database"
	"api/database/models"
	gameHelpers "api/lib/game"
	"encoding/json"
	"errors"

	uuid "github.com/satori/go.uuid"
)

type PlayerState struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	IsConnected bool      `json:"isConnected"`
}

func isConnected(playerId uuid.UUID, lobby map[*Client]bool) bool {
	for client := range lobby {
		if client.PlayerId == playerId {
			return true
		}
	}

	return false
}

func createPlayersUpdate(gameId uuid.UUID, pool *ConnectionPool) ([]PlayerState, error) {
	update := []PlayerState{}
	var game models.Game

	database := database.Get()
	playersError := database.Preload("Players").Where("id = ?", gameId).First(&game).Error
	if playersError != nil {
		return update, errors.New("could not find players")
	}

	lobby := pool.Lobbies[gameId]
	for _, player := range game.Players {
		update = append(update, PlayerState{
			Id:          player.Id,
			Name:        player.Name,
			IsConnected: isConnected(player.Id, lobby),
		})
	}

	return update, nil
}

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

	playersUpdate, playersError := createPlayersUpdate(client.GameId, pool)
	if playersError != nil {
		return errors.New("could not find players")
	}

	// Convert to JSON.
	updateJson, jsonError := json.Marshal(&playersUpdate)
	if jsonError != nil {
		return errors.New("could not read player names")
	}

	// Broadcast the joined message to all plaers in the game.
	pool.Broadcast <- ServerMessage{
		Type:    ServerMessageTypeJoined,
		Content: string(updateJson),
		GameId:  client.GameId,
	}

	return nil
}
