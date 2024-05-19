package websocket

import (
	"api/database"
	"api/database/models"
	gameHelpers "api/lib/game"
	"encoding/json"
	"errors"

	uuid "github.com/satori/go.uuid"
)

func includesPlayer(players []models.Player, playerId uuid.UUID) bool {
	for _, player := range players {
		if player.Id == playerId {
			return true
		}
	}

	return false
}

func findPlayerClient(playerId uuid.UUID, lobby map[*Client]bool) *Client {
	for client := range lobby {
		if client.PlayerId == playerId {
			return client
		}
	}

	return nil
}

type KickPlayerPayload struct {
	PlayerId uuid.UUID `json:"playerId"`
}

func KickPlayerHandler(message ClientMessage, client *Client, pool *ConnectionPool) error {
	var game models.Game

	database := database.Get()
	gameError := database.Preload("Owner").Preload("Players").Where("id = ?", client.GameId).First(&game).Error
	if gameError != nil {
		return errors.New("could not find game")
	}

	// Check if the player is the owner of the game.
	if game.Owner.Id != client.PlayerId {
		return errors.New("only the owner can kick players")
	}

	playerIdToKick, playerToKickUuidError := uuid.FromString(message.Payload)
	if playerToKickUuidError != nil {
		return errors.New("player to kick could not be identified")
	}

	// Check if the player is in the game.
	playerInGame := includesPlayer(game.Players, playerIdToKick)
	if playerInGame {
		// Remove the player from the game.
		removeError := gameHelpers.RemovePlayerFromGame(client.GameId, playerIdToKick)
		if removeError != nil {
			return errors.New("could not kick player")
		}
	}

	// Remove the player from the lobby.
	clientToKick := findPlayerClient(playerIdToKick, pool.Lobbies[client.GameId])
	if clientToKick != nil {
		delete(pool.Lobbies[client.GameId], clientToKick)
	}

	// Broadcast the updated player list.
	broadcastError := BroadcastPlayersUpdate(client, pool)
	if broadcastError != nil {
		return errors.New("could not broadcast player list")
	}

	// Convert to JSON.
	kickPayload := KickPlayerPayload{PlayerId: playerIdToKick}
	kickPayloadJson, jsonError := json.Marshal(&kickPayload)
	if jsonError != nil {
		return errors.New("could not read player names")
	}

	// Let the kicked player know he was kicked.
	clientToKick.Notify(ServerMessage{
		Type:    ServerMessageTypeKickedPlayer,
		Payload: string(kickPayloadJson),
	})

	return nil
}
