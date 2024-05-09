package websocket

import (
	"api/database/models"

	gorilla "github.com/gorilla/websocket"
)

type ConnectionState struct {
	Connection *gorilla.Conn
	Game       *models.Game
	Player     *models.Player
}
