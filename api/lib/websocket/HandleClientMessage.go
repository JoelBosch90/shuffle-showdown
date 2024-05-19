package websocket

import (
	uuid "github.com/satori/go.uuid"
)

type PlayerNames struct {
	Players []string `json:"players"`
}

func HandleClientMessage(message ClientMessage, client *Client, pool *ConnectionPool) {
	// Players should identify themselves with every message.
	playerId, uuidError := uuid.FromString(message.PlayerId.String())
	if uuidError != nil || client.PlayerId != playerId {
		client.SendError("player identification failed")
		return
	}

	switch message.Type {
	case ClientMessageTypeJoin:

		joinError := JoinHandler(message, client, pool)
		if joinError != nil {
			client.SendError(joinError.Error())
			return
		}

	case ClientMessageTypeKickPlayer:

		kickError := KickPlayerHandler(message, client, pool)
		if kickError != nil {
			client.SendError(kickError.Error())
			return
		}

	default:
		// TODO: think of a better default option.
		BroadcastPlayersUpdate(client, pool)
	}
}
