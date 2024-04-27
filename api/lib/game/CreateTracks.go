package game

import (
	"api/database/models"
	databaseHelpers "api/lib/database_helpers"
	spotifyModels "api/lib/spotify_models"

	"github.com/jinzhu/gorm"
)

func CreateTracks(database *gorm.DB, tracks []spotifyModels.Track) error {
	var artistsToCreate []spotifyModels.Artist
	var tracksToCreate []interface{}

	for _, track := range tracks {
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
		return createArtistError
	}

	// Try to update the existing record
	upsertError := databaseHelpers.Upsert(database, tracksToCreate)
	if upsertError != nil {
		return upsertError
	}

	return nil
}
