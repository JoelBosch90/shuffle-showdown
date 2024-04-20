package v1

import (
	"api/lib/spotify"
	"log"
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

func WebSocket(context *gin.Context) {
	connection, error := upgrader.Upgrade(context.Writer, context.Request, nil)
	if error != nil {
		return
	}

	defer connection.Close()
	connection.WriteMessage(websocket.TextMessage, []byte("Hello, client!"))

	token, tokenError := spotify.GetAccessToken()
	if tokenError != nil {
		log.Println(tokenError)
		return
	}
	log.Println(token)

	// Simply echo messages for now.
	for {
		messageType, message, error := connection.ReadMessage()
		if error != nil {
			break
		}
		connection.WriteMessage(messageType, message)
	}
}
