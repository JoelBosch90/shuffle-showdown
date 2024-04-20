package game

import (
	"api/database/models"
	"api/lib/spotify"
	"log"

	"github.com/jinzhu/gorm"
)

func CreateGame(info spotify.PlayListInfo, database *gorm.DB) (models.Game, error) {
	var lastSongAdded string = ""
	game := models.Game{PlayListId: info.Id}

	log.Printf("Number of tracks included: %d", len(info.Tracks.Items))
	log.Printf("Tracks limit: %d", info.Tracks.Limit)
	log.Printf("Tracks total: %d", info.Tracks.Total)

	transactionError := database.Transaction(func(transaction *gorm.DB) error {
		// TODO: Clean up nesting
		// TODO: Test
		// TODO: Fetch more tracks if list > 100?
		for _, item := range info.Tracks.Items {
			if item.AddedAt > lastSongAdded {
				lastSongAdded = item.AddedAt
			}

			releaseYear, releaseMonth, releaseDay := ConvertReleaseDateToIntegers(item.Track.Album.ReleaseDate)
			createTrackError := database.Create(&models.Track{
				ID:           item.Track.Id,
				Name:         item.Track.Name,
				ReleaseYear:  releaseYear,
				ReleaseMonth: releaseMonth,
				ReleaseDay:   releaseDay,
			}).Error
			if createTrackError != nil {
				return createTrackError
			}

			for _, artist := range item.Artists {
				createArtistError := database.Create(&models.Artist{
					ID:   artist.Id,
					Name: artist.Name,
				}).Error
				if createArtistError != nil {
					return createArtistError
				}
			}
		}

		createPlayListError := database.Create(&models.PlayList{
			ID:            info.Id,
			Name:          info.Name,
			LastSongAdded: lastSongAdded,
			TracksTotal:   uint(info.Tracks.Total),
		}).Error
		if createPlayListError != nil {
			return createPlayListError
		}

		createGameError := database.Create(&game).Error
		if createGameError != nil {
			return createGameError
		}

		return nil
	})
	if transactionError != nil {
		return models.Game{}, transactionError
	}

	return game, nil
}
