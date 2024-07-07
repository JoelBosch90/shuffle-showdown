package websocket

import (
	gameHelpers "api/lib/game"
	"encoding/json"
	"errors"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
)

type SelectionUpdate struct {
	SentAt   time.Time          `json:"sentAt"`
	PlayerId uuid.UUID          `json:"playerId"`
	Answer   gameHelpers.Answer `json:"answer"`
}

func AnswerSelectionUpdateHandler(message ClientMessage, client *Client, pool *ConnectionPool) error {

	log.Println("ClientMessageUpdateSelection", message)

	var answer gameHelpers.Answer
	answerParseError := json.Unmarshal([]byte(message.Payload), &answer)
	if answerParseError != nil {
		return errors.New("could not parse answer")
	}

	pool.Broadcast <- ServerMessage{
		Type: ServerMessageTypeAnswerSelectionUpdate,
		Payload: SelectionUpdate{
			SentAt:   time.Now(),
			PlayerId: client.PlayerId,
			Answer:   answer,
		},
		GameId: client.GameId,
	}

	return nil
}
