package game

import (
	"api/database/models"
	spotifyModels "api/lib/spotify_models"
	"log"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func CreateGame(info spotifyModels.Playlist, database *gorm.DB) (models.Game, error) {
	game := models.Game{
		Id:         uuid.NewV4(),
		PlaylistId: info.Id,
	}
	log.Println("Creating game with playlist ID", info.Id)

	playlistError := CreatePlaylist(info, database)
	if playlistError != nil {
		return models.Game{}, playlistError
	}

	createGameError := database.Create(&game).Error
	if createGameError != nil {
		return models.Game{}, createGameError
	}

	return game, nil
}
