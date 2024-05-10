package websocket

import (
	uuid "github.com/satori/go.uuid"
)

func HandleClientMessage(message ClientMessage, client *Client, pool *ConnectionPool) {
	switch message.Type {
	case string(ClientMessageTypeJoin):

		// Check if the player has already identified himself.
		playerId, uuidError := uuid.FromString(message.PlayerId.String())
		if uuidError != nil || client.Player.Id != playerId {
			client.Notify(ServerMessage{
				Type:    ServerMessageTypeError,
				Content: "Player identification failed",
			})
			return
		}

		pool.Broadcast <- ServerMessage{
			Type:    ServerMessageTypeJoined,
			Content: "Welcome " + client.Player.Name + "!",
			Game:    client.Game,
		}

	default:
		pool.Broadcast <- ServerMessage{
			Type:    ServerMessageTypeJoined,
			Content: "Hello client!",
			Game:    client.Game,
		}
	}
}
