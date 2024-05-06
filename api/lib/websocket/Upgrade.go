package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const BUFFER_SIZE = 1024

var upgrader = websocket.Upgrader{
	ReadBufferSize:  BUFFER_SIZE,
	WriteBufferSize: BUFFER_SIZE,
	CheckOrigin: func(request *http.Request) bool {
		return true
	},
}

func Upgrade(context *gin.Context) (*websocket.Conn, error) {

	// Upgrade the connection to a Websocket connection.
	connection, error := upgrader.Upgrade(context.Writer, context.Request, nil)
	if error != nil {
		return nil, error
	}

	return connection, nil
}
