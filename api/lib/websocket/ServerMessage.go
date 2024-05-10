package websocket

import "api/database/models"

type ServerMessageType string

const (
	ServerMessageTypeJoined ServerMessageType = "joined"
	ServerMessageTypeError  ServerMessageType = "error"
)

type ServerMessage struct {
	Type    ServerMessageType `json:"type"`
	Content string            `json:"content"`
	Game    *models.Game      `json:"-"`
}
