package game

import (
	"api/database/models"
	databaseHelpers "api/lib/database_helpers"
	spotifyModels "api/lib/spotify_models"

	"github.com/jinzhu/gorm"
)

func CreateArtist(database *gorm.DB, artist spotifyModels.Artist) error {
	upsertError := databaseHelpers.Upsert(database, &models.Artist{
		ID:   artist.Id,
		Name: artist.Name,
	})
	if upsertError != nil {
		return upsertError
	}

	return nil
}
