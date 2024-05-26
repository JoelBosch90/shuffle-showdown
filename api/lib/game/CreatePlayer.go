package game

import (
	"api/database/models"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func CreatePlayer(name string, database *gorm.DB) (*models.Player, error) {
	player := models.Player{
		Id:     uuid.NewV4(),
		Secret: uuid.NewV4(),
		Name:   name,
	}

	playerError := database.Create(&player).Error
	if playerError != nil {
		return nil, playerError
	}

	return &player, nil
}
