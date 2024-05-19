package websocket

import uuid "github.com/satori/go.uuid"

type ServerMessageType string

const (
	ServerMessageTypeGameUpdate   ServerMessageType = "game-update"
	ServerMessageTypeKickedPlayer ServerMessageType = "kicked-player"
	ServerMessageTypeError        ServerMessageType = "error"
)

type ServerMessage struct {
	Type    ServerMessageType `json:"type"`
	Payload interface{}       `json:"payload"`
	GameId  uuid.UUID         `json:"-"`
}
