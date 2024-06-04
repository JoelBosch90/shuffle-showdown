package websocket

import (
	"api/database"
	"api/database/models"
	gameHelpers "api/lib/game"
	"errors"

	uuid "github.com/satori/go.uuid"
)

func findPlayerClient(playerId uuid.UUID, lobby map[*Client]struct{}) *Client {
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

	clientToKick.Notify(ServerMessage{
		Type:    ServerMessageTypeKickedPlayer,
		Payload: KickPlayerPayload{PlayerId: playerIdToKick},
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

	broadcastError := BroadcastGameUpdate(client, pool)
	if broadcastError != nil {
		return errors.New("could not broadcast game update")
	}

	return nil
}
