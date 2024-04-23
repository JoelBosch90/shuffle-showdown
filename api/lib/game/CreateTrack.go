package game

import (
	"api/database/models"
	databaseHelpers "api/lib/database_helpers"
	spotifyModels "api/lib/spotify_models"

	"github.com/jinzhu/gorm"
)

func CreateTrack(database *gorm.DB, track spotifyModels.Track) error {
	for _, artist := range track.Artists {
		createArtistError := CreateArtist(database, artist)
		if createArtistError != nil {
			return createArtistError
		}
	}

	// Try to update the existing record
	releaseYear, releaseMonth, releaseDay := ConvertReleaseDateToIntegers(track.Album.ReleaseDate)
	upsertError := databaseHelpers.Upsert(database, &models.Track{
		ID:           track.Id,
		Name:         track.Name,
		ReleaseYear:  releaseYear,
		ReleaseMonth: releaseMonth,
		ReleaseDay:   releaseDay,
	})
	if upsertError != nil {
		return upsertError
	}

	return nil
}
