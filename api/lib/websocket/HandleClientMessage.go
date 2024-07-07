package websocket

import (
	uuid "github.com/satori/go.uuid"
)

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
		}

	case ClientMessageTypeKickPlayer:

		kickError := KickPlayerHandler(message, client, pool)
		if kickError != nil {
			client.SendError(kickError.Error())
		}

	case ClientMessageStartGame:

		startGameError := StartGameHandler(message, client, pool)
		if startGameError != nil {
			client.SendError(startGameError.Error())
		}

	case ClientMessageSubmitAnswer:

		submitAnswerError := SubmitAnswerHandler(message, client, pool)
		if submitAnswerError != nil {
			client.SendError(submitAnswerError.Error())
		}

	case ClientMessageUpdateRequest:

		broadcastError := BroadcastGameUpdate(client, pool)
		if broadcastError != nil {
			client.SendError(broadcastError.Error())
		}

	case ClientMessageUpdateAnswerSelection:

		answerSelectionUpdateError := AnswerSelectionUpdateHandler(message, client, pool)
		if answerSelectionUpdateError != nil {
			client.SendError(answerSelectionUpdateError.Error())
		}

	default:
		BroadcastGameUpdate(client, pool)
	}
}
