package spotify

import (
	databaseHelpers "api/database"
	"api/database/models"
	spotifyModels "api/lib/spotify/models"
	"errors"

	"github.com/jinzhu/gorm"
)

func assertAlbums(upsertedAlbums []interface{}) ([]models.Album, error) {
	albums := make([]models.Album, len(upsertedAlbums))
	for index, upsertedAlbum := range upsertedAlbums {
		album, ok := upsertedAlbum.(models.Album)
		if !ok {
			return []models.Album{}, errors.New("could not convert upsertedAlbum to []models.Album")
		}
		albums[index] = album
	}

	return albums, nil
}

func CreateAlbums(database *gorm.DB, items []spotifyModels.Item) ([]models.Album, error) {
	var albumsToCreate []interface{}

	for _, item := range items {
		album := item.Album

		releaseYear, releaseMonth, releaseDay := ConvertReleaseDateToIntegers(album.ReleaseDate)

		if releaseYear == 0 {
			continue
		}

		albumsToCreate = append(albumsToCreate, &models.Album{
			Id:           album.Id,
			Name:         album.Name,
			ReleaseYear:  releaseYear,
			ReleaseMonth: releaseMonth,
			ReleaseDay:   releaseDay,
		})
	}

	upsertedAlbums, upsertError := databaseHelpers.Upsert(database, albumsToCreate)
	if upsertError != nil {
		return []models.Album{}, upsertError
	}

	createdAlbums, assertError := assertAlbums(upsertedAlbums)
	if assertError != nil {
		return []models.Album{}, errors.New("could not convert upsertedAlbums to []models.Album")
	}

	return createdAlbums, nil
}
