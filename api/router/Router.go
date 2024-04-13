package router

import (
	"api/router/endpoints"
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

func Run() {
	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("game", endpoints.PostGame)
		apiV1.GET("game/:id", endpoints.GetGame)
		apiV1.PATCH("game/:id", endpoints.PatchGame)
		apiV1.DELETE("game/:id", endpoints.DeleteGame)
		apiV1.GET("ws", func(context *gin.Context) {
			connection, error := upgrader.Upgrade(context.Writer, context.Request, nil)
			if error != nil {
				return
			}

			defer connection.Close()
			connection.WriteMessage(websocket.TextMessage, []byte("Hello, client!"))
		})
	}

	router.Run()
}
