package router

import (
	"api/router/endpoints"

	"github.com/gin-gonic/gin"
)

func Run() {
	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("game", endpoints.PostGame)
		apiV1.GET("game/:id", endpoints.GetGame)
		apiV1.PATCH("game/:id", endpoints.PatchGame)
		apiV1.DELETE("game/:id", endpoints.DeleteGame)
	}

	router.Run()
}
