package router

import (
	game "api/router/endpoints/v1/game"
	player "api/router/endpoints/v1/player"

	"github.com/gin-gonic/gin"
)

func Run() {
	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("game", game.Post)
		apiV1.GET("game/:id", game.Get)
		apiV1.PATCH("game/:id", game.Patch)
		apiV1.DELETE("game/:id", game.Delete)

		apiV1.POST("player", player.Post)
		apiV1.PATCH("player/:id", player.Patch)

		apiV1.GET("ws/:id", game.WebSocket)
	}

	router.Run()
}
