package game

import (
	"api/database/models"
	databaseHelpers "api/lib/database_helpers"
	spotifyModels "api/lib/spotify_models"
	"errors"

	"github.com/jinzhu/gorm"
)

func assertArtists(upsertedArtists []interface{}) ([]models.Artist, error) {
	artists := make([]models.Artist, len(upsertedArtists))
	for index, upsertedArtist := range upsertedArtists {
		artist, ok := upsertedArtist.(models.Artist)
		if !ok {
			return []models.Artist{}, errors.New("could not convert upsertedArtists to []models.Artist")
		}
		artists[index] = artist
	}

	return artists, nil
}

func CreateArtists(database *gorm.DB, artists []spotifyModels.Artist) ([]models.Artist, error) {
	var artistsToCreate []interface{}
	for _, artist := range artists {
		artistsToCreate = append(artistsToCreate, &models.Artist{
			Id:   artist.Id,
			Name: artist.Name,
		})
	}

	upsertedArtists, upsertError := databaseHelpers.Upsert(database, artistsToCreate)
	if upsertError != nil {
		return []models.Artist{}, upsertError
	}

	createdArtists, assertError := assertArtists(upsertedArtists)
	if assertError != nil {
		return []models.Artist{}, errors.New("could not convert upsertedTracks to []models.Track")
	}

	return createdArtists, nil
}
