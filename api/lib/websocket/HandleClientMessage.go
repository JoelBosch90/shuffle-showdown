package websocket

import (
	uuid "github.com/satori/go.uuid"
)

func HandleClientMessage(message ClientMessage, client *Client, pool *ConnectionPool) {
	switch message.Type {
	case "join":

		// Check if the player has already identified himself.
		playerId, uuidError := uuid.FromString(message.PlayerId.String())
		if uuidError != nil || client.Player.Id != playerId {
			client.Notify(ServerMessage{
				Type:    "error",
				Content: "Player identification failed",
			})
			return
		}

		pool.Broadcast <- ServerMessage{
			Type:    "joined",
			Content: "Welcome " + client.Player.Name + "!",
			Game:    client.Game,
		}

	default:
		pool.Broadcast <- ServerMessage{
			Type:    "joined",
			Content: "Hello client!",
			Game:    client.Game,
		}
	}
}
