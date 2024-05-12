package websocket

import uuid "github.com/satori/go.uuid"

type ServerMessageType string

const (
	ServerMessageTypePlayersUpdate ServerMessageType = "players-update"
	ServerMessageTypeError         ServerMessageType = "error"
)

type ServerMessage struct {
	Type    ServerMessageType `json:"type"`
	Content string            `json:"content"`
	GameId  uuid.UUID         `json:"-"`
}
