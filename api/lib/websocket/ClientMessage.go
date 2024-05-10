package websocket

import uuid "github.com/satori/go.uuid"

type ClientMessageType string

const (
	ClientMessageTypeJoin ClientMessageType = "join"
)

type ClientMessage struct {
	Type     string
	Content  string
	PlayerId uuid.UUID
}
