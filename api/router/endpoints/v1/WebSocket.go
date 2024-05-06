package v1

import (
	"api/database"
	"api/database/models"
	websocketHelpers "api/lib/websocket"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var connnectionPool = make(map[*websocket.Conn]*PlayerState)

type PlayerState struct {
	Conn   *websocket.Conn
	Game   *models.Game
	Player *models.Player
}

type ServerMessage struct {
	Type    string
	Content string
}

type ClientMessage struct {
	Type     string
	Content  string
	PlayerId uuid.UUID
}

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
	connection, error := websocketHelpers.Upgrade(context)
	if error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not upgrade to Websocket connection"})
		return
	}

	defer connection.Close()

	// Add the connection to the connection pool.
	connnectionPool[connection] = &PlayerState{
		Conn: connection,
		Game: &game,
	}

	connection.WriteMessage(websocket.TextMessage, []byte("Hello, client!"))

	// Simply echo messages for now.
	for {
		var message ClientMessage
		messageError := connection.ReadJSON(&message)
		if messageError != nil {
			connection.WriteJSON(ServerMessage{
				Type:    "error",
				Content: "Error reading message",
			})
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
				connection.WriteJSON(ServerMessage{
					Type:    "error",
					Content: "Error identifying player",
				})
				continue
			}

			// Store the player in the player state.
			playerState.Player = &models.Player{}
			databaseError := database.Where("id = ? AND secret = ?", playerId, secret).First(playerState.Player).Error
			if databaseError != nil {
				connection.WriteJSON(ServerMessage{
					Type:    "error",
					Content: "Unknown player",
				})
			}

			log.Println("Player identified:", playerState.Player)
		default:
			continue
		}

		// log.Println(messageType, message)

		// connection.WriteMessage(messageType, message)
	}
}
