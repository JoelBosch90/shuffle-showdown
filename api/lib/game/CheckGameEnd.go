package game

import (
	"api/database"
	"api/database/models"
	"log"

	uuid "github.com/satori/go.uuid"
)

func hasFinished(game *models.Game) bool {
	log.Println("CHECKING....", game.SongsToWin)

	for _, player := range game.Players {
		wonTracksInGame := FilterWonTracksByGameId(game.Id, player.WonTracks)
		log.Println("CHECKING....", wonTracksInGame, uint(len(wonTracksInGame)), game.SongsToWin)

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
		log.Println("GAME HAS NOT FINISHED", game.Id)
		return false
	}

	log.Println("GAME HAS FINISHED", game.Id)

	database.Model(&models.Game{}).Where("id = ?", gameId).Updates(models.Game{HasFinished: true})

	return true
}
