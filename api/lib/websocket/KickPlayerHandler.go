package websocket

import (
	"api/database"
	"api/database/models"
	gameHelpers "api/lib/game"
	"encoding/json"
	"errors"

	uuid "github.com/satori/go.uuid"
)

func findPlayerClient(playerId uuid.UUID, lobby map[*Client]bool) *Client {
	for client := range lobby {
		if client.PlayerId == playerId {
			return client
		}
	}

	return nil
}

func kickPlayerFromConnectionPool(playerIdToKick uuid.UUID, client *Client, pool *ConnectionPool) error {
	clientToKick := findPlayerClient(playerIdToKick, pool.Lobbies[client.GameId])
	if clientToKick == nil {
		return nil
	}

	kickPayload := KickPlayerPayload{PlayerId: playerIdToKick}
	kickPayloadJson, jsonError := json.Marshal(&kickPayload)
	if jsonError != nil {
		return errors.New("could not read player names")
	}

	clientToKick.Notify(ServerMessage{
		Type:    ServerMessageTypeKickedPlayer,
		Payload: string(kickPayloadJson),
	})

	delete(pool.Lobbies[client.GameId], clientToKick)
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

	if game.Owner.Id != client.PlayerId {
		return errors.New("only the owner can kick players")
	}

	playerIdToKick, playerToKickUuidError := uuid.FromString(message.Payload)
	if playerToKickUuidError != nil {
		return errors.New("player to kick could not be identified")
	}

	banError := gameHelpers.BanPlayerFromGame(client.GameId, playerIdToKick)
	if banError != nil {
		return errors.New("could not ban player")
	}

	kickFromClientError := kickPlayerFromConnectionPool(playerIdToKick, client, pool)
	if kickFromClientError != nil {
		return errors.New("could not kick player")
	}

	broadcastError := BroadcastPlayersUpdate(client, pool)
	if broadcastError != nil {
		return errors.New("could not broadcast player list")
	}

	return nil
}
