package game

import (
	"api/database"
	"api/database/models"
	"api/lib/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WebSocket(context *gin.Context) {
	gameId := context.Param("id")
	database := database.Get()
	var game models.Game

	databaseError := database.Where("id = ?", gameId).First(&game).Error
	if databaseError != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	connection, error := websocket.Upgrade(context)
	if error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not upgrade to Websocket connection"})
		return
	}

	secret, secretError := context.Cookie("playerSecret")
	if secretError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unknown player"})
		return

	}

	player := &models.Player{}
	playerError := database.Where("secret = ?", secret).First(player).Error
	if playerError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unknown player"})
		return
	}

	connectionPool := websocket.GetConnectionPool()
	client := &websocket.Client{
		Pool:             connectionPool,
		Connection:       connection,
		OutgoingMessages: make(chan websocket.ServerMessage, 256),
		GameId:           game.Id,
		PlayerId:         player.Id,
	}
	connectionPool.Register <- client

	go client.Read()
	go client.Write()
}
