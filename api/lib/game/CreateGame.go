package game

import (
	"api/database/models"
	databaseHelpers "api/lib/database_helpers"
	spotifyModels "api/lib/spotify_models"

	"github.com/jinzhu/gorm"
)

func CreateGame(info spotifyModels.PlayList, database *gorm.DB) (models.Game, error) {
	var lastSongAdded string = ""
	game := models.Game{PlayListId: info.Id}

	for _, item := range info.Tracks.Items {
		if item.AddedAt > lastSongAdded {
			lastSongAdded = item.AddedAt
		}

		createTrackError := CreateTrack(database, item.Track)
		if createTrackError != nil {
			return models.Game{}, createTrackError
		}
	}

	upsertPlayListError := databaseHelpers.Upsert(database, &models.PlayList{
		ID:            info.Id,
		Name:          info.Name,
		LastSongAdded: lastSongAdded,
		TracksTotal:   uint(info.Tracks.Total),
	})
	if upsertPlayListError != nil {
		return models.Game{}, upsertPlayListError
	}

	createGameError := database.Create(&game).Error
	if createGameError != nil {
		return models.Game{}, createGameError
	}

	return game, nil
}
