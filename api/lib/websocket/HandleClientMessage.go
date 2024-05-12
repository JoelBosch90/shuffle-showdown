package websocket

import (
	"api/lib/game"
	"encoding/json"

	uuid "github.com/satori/go.uuid"
)

type PlayerNames struct {
	Players []string `json:"players"`
}

func HandleClientMessage(message ClientMessage, client *Client, pool *ConnectionPool) {
	switch message.Type {
	case ClientMessageTypeJoin:

		// Check if the player has already identified himself.
		playerId, uuidError := uuid.FromString(message.PlayerId.String())
		if uuidError != nil || client.PlayerId != playerId {
			client.Notify(ServerMessage{
				Type:    ServerMessageTypeError,
				Content: "Player identification failed",
			})
			return
		}

		names, namesError := game.GetPlayerNames(client.GameId)
		if namesError != nil {
			client.Notify(ServerMessage{
				Type:    ServerMessageTypeError,
				Content: "Could not read player names",
			})
			return
		}

		namesJson, jsonError := json.Marshal(&PlayerNames{
			Players: names,
		})
		if jsonError != nil {
			client.Notify(ServerMessage{
				Type:    ServerMessageTypeError,
				Content: "Could not read player names",
			})
			return
		}

		pool.Broadcast <- ServerMessage{
			Type:    ServerMessageTypeJoined,
			Content: string(namesJson),
			GameId:  client.GameId,
		}

	default:
		pool.Broadcast <- ServerMessage{
			Type:    ServerMessageTypeJoined,
			Content: "Hello client!",
			GameId:  client.GameId,
		}
	}
}
