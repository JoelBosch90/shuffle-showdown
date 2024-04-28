package game

import (
	"api/database/models"
	spotifyModels "api/lib/spotify_models"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func CreateGame(info spotifyModels.Playlist, database *gorm.DB) (models.Game, error) {
	playlist, playlistError := CreatePlaylist(info, database)
	if playlistError != nil {
		return models.Game{}, playlistError
	}

	game := models.Game{
		Id:       uuid.NewV4(),
		Playlist: playlist,
	}

	createGameError := database.Create(&game).Error
	if createGameError != nil {
		return models.Game{}, createGameError
	}

	return game, nil
}
