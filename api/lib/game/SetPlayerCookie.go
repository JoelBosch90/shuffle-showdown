package game

import (
	"api/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetPlayerCookie(context *gin.Context, player models.Player) {
	cookie := http.Cookie{
		Name:     "playerSecret",
		Value:    player.Secret.String(),
		HttpOnly: true,
	}
	http.SetCookie(context.Writer, &cookie)
}
