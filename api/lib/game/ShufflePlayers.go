package game

import (
	database "api/database"
	"api/database/models"
	"errors"
	"log"
	"math/rand"

	uuid "github.com/satori/go.uuid"
)

func ShufflePlayers(gameId uuid.UUID) error {
	database := database.Get()
	var game models.Game

	loadGameError := database.Preload("Players").Where("id = ?", gameId).First(&game).Error
	if loadGameError != nil {
		return errors.New("could not load game")
	}

	players := game.Players
	rand.Shuffle(len(players), func(i, j int) { players[i], players[j] = players[j], players[i] })

	for order, player := range players {
		log.Println("ORDER: ", order)
		log.Println("PLAYER: ", player)
		database.Model(&models.GamePlayer{}).Where("game_id = ? AND player_id = ?", gameId, player.Id).Updates(models.GamePlayer{Order: uint(order)})
	}

	return nil
}
