package game

import (
	"api/database/models"
	databaseHelpers "api/lib/database_helpers"
	spotifyModels "api/lib/spotify_models"

	"github.com/jinzhu/gorm"
)

func CreateArtists(database *gorm.DB, artists []spotifyModels.Artist) error {
	var artistsToCreate []interface{}
	for _, artist := range artists {
		artistsToCreate = append(artistsToCreate, &models.Artist{
			Id:   artist.Id,
			Name: artist.Name,
		})
	}

	upsertError := databaseHelpers.Upsert(database, artistsToCreate)
	if upsertError != nil {
		return upsertError
	}

	return nil
}
