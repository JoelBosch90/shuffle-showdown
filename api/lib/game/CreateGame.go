package game

import (
	"api/database/models"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func CreateGame(playlistId string, countryCode string, player models.Player, database *gorm.DB) (models.Game, error) {
	game := models.Game{
		Id:          uuid.NewV4(),
		CountryCode: countryCode,
		PlaylistId:  playlistId,
		SongsToWin:  10,
		Owner:       player,
		Players:     []models.Player{player},
	}

	createGameError := database.Create(&game).Error
	if createGameError != nil {
		return models.Game{}, createGameError
	}

	return game, nil
}
