package game

import (
	"api/database/models"
	spotify "api/lib/spotify"
	spotifyModels "api/lib/spotify/models"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func CreateGame(info spotifyModels.Playlist, player models.Player, database *gorm.DB) (models.Game, error) {
	playlist, playlistError := spotify.CreatePlaylist(info, database)
	if playlistError != nil {
		return models.Game{}, playlistError
	}

	game := models.Game{
		Id:       uuid.NewV4(),
		Playlist: playlist,
		Owner:    player,
		Players:  []models.Player{player},
	}

	createGameError := database.Create(&game).Error
	if createGameError != nil {
		return models.Game{}, createGameError
	}

	return game, nil
}
