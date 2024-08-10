package router

import (
	"api/database"
	"api/database/models"
	"api/lib/verification"
	game "api/router/endpoints/v1/game"
	mostPlayedPlaylists "api/router/endpoints/v1/mostPlayedPlaylists"
	player "api/router/endpoints/v1/player"

	"github.com/gin-gonic/gin"
)

func testVerification(context *gin.Context) {
	database := database.Get()
	trackTitle := "Sapés comme jamais (feat. Niska) - Pilule bleue"

	var tracks []models.Track
	database.Model(&models.Track{}).Where("name = ?", trackTitle).Scan(&tracks)

	for _, track := range tracks {
		verification.VerifyTrack(database, track.Id)
	}
}

func Run() {
	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("playlists/most-played", mostPlayedPlaylists.Get)

		apiV1.POST("game", game.Post)
		apiV1.GET("game/:id", game.Get)
		apiV1.PATCH("game/:id", game.Patch)
		apiV1.DELETE("game/:id", game.Delete)

		apiV1.POST("player", player.Post)
		apiV1.PATCH("player/:id", player.Patch)

		apiV1.GET("ws/:id", game.WebSocket)
		apiV1.GET("test-verification", testVerification)
	}

	router.Run()
}
