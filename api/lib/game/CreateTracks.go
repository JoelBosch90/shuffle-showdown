package game

import (
	"api/database/models"
	databaseHelpers "api/lib/database_helpers"
	spotifyModels "api/lib/spotify_models"

	"github.com/jinzhu/gorm"
)

func CreateTracks(database *gorm.DB, items []spotifyModels.Item) (string, error) {
	var artistsToCreate []spotifyModels.Artist
	var tracksToCreate []interface{}
	var lastSongAdded string = ""

	for _, item := range items {
		lastSongAdded = item.AddedAt
		track := item.Track
		artistsToCreate = append(artistsToCreate, track.Artists...)

		releaseYear, releaseMonth, releaseDay := ConvertReleaseDateToIntegers(track.Album.ReleaseDate)
		tracksToCreate = append(tracksToCreate, &models.Track{
			Id:           track.Id,
			Name:         track.Name,
			ReleaseYear:  releaseYear,
			ReleaseMonth: releaseMonth,
			ReleaseDay:   releaseDay,
		})
	}

	createArtistError := CreateArtists(database, artistsToCreate)
	if createArtistError != nil {
		return lastSongAdded, createArtistError
	}

	// Try to update the existing record
	upsertError := databaseHelpers.Upsert(database, tracksToCreate)
	if upsertError != nil {
		return lastSongAdded, upsertError
	}

	return lastSongAdded, nil
}
