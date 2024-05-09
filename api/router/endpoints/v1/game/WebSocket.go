package game

import (
	"api/database"
	"api/database/models"
	"api/lib/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

var connnectionPool = make(map[*gorilla.Conn]*websocket.ConnectionState)

func WebSocket(context *gin.Context) {
	id := context.Param("id")
	database := database.Get()
	var game models.Game

	// Get the game from the database.
	databaseError := database.Where("id = ?", id).First(&game).Error
	if databaseError != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	// Upgrade the connection to a Websocket connection.
	connection, error := websocket.Upgrade(context)
	if error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not upgrade to Websocket connection"})
		return
	}

	defer connection.Close()

	// Add the connection to the connection pool.
	connnectionPool[connection] = &websocket.ConnectionState{
		Connection: connection,
		Game:       &game,
	}

	// Listen for messages.
	for {
		var message websocket.ClientMessage
		messageError := connection.ReadJSON(&message)
		if messageError != nil {
			connection.WriteJSON(websocket.ServerMessage{
				Type:    "error",
				Content: "Error reading message",
			})
			connection.Close()
			continue
		}

		switch message.Type {
		case "join":
			// Identify the player.
			playerState := connnectionPool[connection]

			// Check if the player has already identified himself.
			playerId := message.PlayerId.String()
			secret, secretError := context.Cookie("playerSecret")
			if secretError != nil {
				connection.WriteJSON(websocket.ServerMessage{
					Type:    "error",
					Content: "Error identifying player",
				})
				continue
			}

			// Store the player in the player state.
			playerState.Player = &models.Player{}
			databaseError := database.Where("id = ? AND secret = ?", playerId, secret).First(playerState.Player).Error
			if databaseError != nil {
				connection.WriteJSON(websocket.ServerMessage{
					Type:    "error",
					Content: "Unknown player",
				})
			}

			connection.WriteMessage(gorilla.TextMessage, []byte("Hello, "+playerState.Player.Name+"!"))
		default:

			connection.WriteMessage(gorilla.TextMessage, []byte("Hello, client!"))
			continue
		}
	}
}
