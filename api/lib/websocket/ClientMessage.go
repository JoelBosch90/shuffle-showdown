package websocket

import uuid "github.com/satori/go.uuid"

type ClientMessageType string

const (
	ClientMessageTypeJoin       ClientMessageType = "join"
	ClientMessageTypeKickPlayer ClientMessageType = "kick-player"
	ClientMessageStartGame      ClientMessageType = "start-game"
	ClientMessageSubmitAnswer   ClientMessageType = "submit-answer"
	ClientMessageUpdateRequest  ClientMessageType = "update-request"
)

type ClientMessage struct {
	Type     ClientMessageType `json:"type"`
	Payload  string            `json:"payload"`
	PlayerId uuid.UUID         `json:"playerId"`
}
