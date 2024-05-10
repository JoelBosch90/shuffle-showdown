package websocket

import "api/database/models"

type ServerMessage struct {
	Type    string       `json:"type"`
	Content string       `json:"content"`
	Game    *models.Game `json:"-"`
}
