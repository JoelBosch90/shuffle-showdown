package websocket

import uuid "github.com/satori/go.uuid"

type ServerMessageType string

const (
	ServerMessageTypeJoined ServerMessageType = "joined"
	ServerMessageTypeError  ServerMessageType = "error"
)

type ServerMessage struct {
	Type    ServerMessageType `json:"type"`
	Content string            `json:"content"`
	GameId  uuid.UUID         `json:"-"`
}
