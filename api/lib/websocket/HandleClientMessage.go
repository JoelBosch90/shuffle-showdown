package websocket

type PlayerNames struct {
	Players []string `json:"players"`
}

func HandleClientMessage(message ClientMessage, client *Client, pool *ConnectionPool) {
	switch message.Type {
	case ClientMessageTypeJoin:

		joinError := JoinHandler(message, client, pool)
		if joinError != nil {
			client.Notify(ServerMessage{
				Type:    ServerMessageTypeError,
				Content: joinError.Error(),
			})
			return
		}

	default:
		// TODO: think of a better default option.
		BroadcastPlayersUpdate(client, pool)
	}
}
