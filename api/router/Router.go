package router

import (
	v1 "api/router/endpoints/v1"

	"github.com/gin-gonic/gin"
)

func Run() {
	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("game", v1.PostGame)
		apiV1.GET("game/:id", v1.GetGame)
		apiV1.PATCH("game/:id", v1.PatchGameSettings)
		apiV1.DELETE("game/:id", v1.DeleteGame)

		apiV1.PATCH("player/:id", v1.PatchPlayerName)

		apiV1.GET("ws/:id", v1.WebSocket)
	}

	router.Run()
}
