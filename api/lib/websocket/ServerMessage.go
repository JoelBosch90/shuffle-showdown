package websocket

import uuid "github.com/satori/go.uuid"

type ServerMessageType string

const (
	ServerMessageTypePlayersUpdate ServerMessageType = "players-update"
	ServerMessageTypeKickedPlayer  ServerMessageType = "kicked-player"
	ServerMessageTypeError         ServerMessageType = "error"
)

type ServerMessage struct {
	Type    ServerMessageType `json:"type"`
	Payload string            `json:"payload"`
	GameId  uuid.UUID         `json:"-"`
}
