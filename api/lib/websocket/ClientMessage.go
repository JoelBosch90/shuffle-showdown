package websocket

import uuid "github.com/satori/go.uuid"

type ClientMessage struct {
	Type     string
	Content  string
	PlayerId uuid.UUID
}
