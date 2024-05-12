package game

import (
	"api/database"
	"api/database/models"
	"errors"

	uuid "github.com/satori/go.uuid"
)

func GetPlayerNames(gameId uuid.UUID) ([]string, error) {
	var names []string
	var freshGame models.Game

	database := database.Get()
	databaseError := database.Preload("Players").Where("id = ?", gameId).First(&freshGame).Error
	if databaseError != nil {
		return []string{}, errors.New("could not fetch players")
	}

	for _, player := range freshGame.Players {
		names = append(names, player.Name)
	}

	return names, nil
}
