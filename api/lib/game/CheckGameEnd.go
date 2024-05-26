package game

import (
	"api/database"
	"api/database/models"

	uuid "github.com/satori/go.uuid"
)

func hasFinished(game *models.Game) bool {
	for _, player := range game.Players {
		wonTracksInGame := FilterWonTracksByGameId(game.Id, player.WonTracks)

		if uint(len(wonTracksInGame)) >= game.SongsToWin {
			return true
		}
	}

	return false
}

func CheckGameEnd(gameId uuid.UUID) bool {
	database := database.Get()
	var game models.Game

	gameError := database.Preload("Players.WonTracks").Where("id = ?", gameId).First(&game).Error
	if gameError != nil {
		return false
	}

	if !hasFinished(&game) {
		return false
	}

	database.Model(&models.Game{}).Where("id = ?", gameId).Updates(models.Game{HasFinished: true})

	return true
}
