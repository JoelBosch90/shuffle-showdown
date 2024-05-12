package websocket

import "log"

type PlayerNames struct {
	Players []string `json:"players"`
}

func HandleClientMessage(message ClientMessage, client *Client, pool *ConnectionPool) {
	log.Println("Handling message", message)

	switch message.Type {

	case ClientMessageTypeJoin:

		joinError := JoinHandler(message, client, pool)
		if joinError != nil {
			client.Notify(ServerMessage{
				Type:    ServerMessageTypeError,
				Payload: joinError.Error(),
			})
			return
		}

	case ClientMessageTypeKickPlayer:

		kickError := KickPlayerHandler(message, client, pool)
		if kickError != nil {
			client.Notify(ServerMessage{
				Type:    ServerMessageTypeError,
				Payload: kickError.Error(),
			})
			return
		}

	default:
		// TODO: think of a better default option.
		BroadcastPlayersUpdate(client, pool)
	}
}
