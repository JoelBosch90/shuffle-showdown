package spotify

import (
	databaseHelpers "api/database"
	"api/database/models"
	spotifyModels "api/lib/spotify/models"
	"errors"

	"github.com/jinzhu/gorm"
)

func findArtistById(artists []models.Artist, id string) *models.Artist {
	for _, artist := range artists {
		if artist.Id == id {
			return &artist
		}
	}

	return nil
}

func containsAlbum(albums []models.Album, id string) bool {
	for _, album := range albums {
		if album.Id == id {
			return true
		}
	}

	return false
}

func constructTracks(items []spotifyModels.Item, createdArtists []models.Artist, createdAlbums []models.Album) ([]interface{}, []interface{}, error) {
	var tracksToCreate []interface{}
	var trackArtistsToCreate []interface{}

	for _, item := range items {
		track := item.Track

		if track.PreviewUrl == "" || !containsAlbum(createdAlbums, track.Album.Id) {
			continue
		}

		artistIds := []string{}
		for _, artist := range track.Artists {
			artistIds = append(artistIds, artist.Id)
		}

		artists := []models.Artist{}
		for _, artistId := range artistIds {
			artist := findArtistById(createdArtists, artistId)
			if artist != nil {
				artists = append(artists, *artist)
			}

			trackArtistsToCreate = append(trackArtistsToCreate, &models.TrackArtist{
				TrackId:  track.Id,
				ArtistId: artistId,
			})
		}

		tracksToCreate = append(tracksToCreate, &models.Track{
			Id:         track.Id,
			Name:       track.Name,
			AlbumId:    track.Album.Id,
			Artists:    artists,
			PreviewUrl: track.PreviewUrl,
			IsPlayable: track.IsPlayable,
		})
	}

	return tracksToCreate, trackArtistsToCreate, nil
}

func assertTracks(upsertedTracks []interface{}) ([]models.Track, error) {
	tracks := make([]models.Track, len(upsertedTracks))

	for index, upsertedTrack := range upsertedTracks {
		upsertedTrack, ok := upsertedTrack.(models.Track)
		if !ok {
			return []models.Track{}, errors.New("could not convert upsertedTracks to []models.Track")
		}
		tracks[index] = upsertedTrack
	}

	return tracks, nil
}

func CreateTracks(database *gorm.DB, items []spotifyModels.Item) ([]models.Track, error) {
	createdArtists, createdArtistError := CreateArtists(database, items)
	if createdArtistError != nil || len(createdArtists) == 0 {
		return []models.Track{}, createdArtistError
	}

	createdAlbums, createdAlbumsError := CreateAlbums(database, items)
	if createdAlbumsError != nil || len(createdAlbums) == 0 {
		return []models.Track{}, createdAlbumsError
	}

	tracksToCreate, trackArtistsToCreate, constructError := constructTracks(items, createdArtists, createdAlbums)
	if constructError != nil || len(tracksToCreate) == 0 {
		return []models.Track{}, constructError
	}

	upsertedTracks, upsertError := databaseHelpers.Upsert(database, tracksToCreate)
	if upsertError != nil || len(upsertedTracks) == 0 {
		return []models.Track{}, upsertError
	}

	_, trackArtistsUpsertError := databaseHelpers.Upsert(database, trackArtistsToCreate)
	if trackArtistsUpsertError != nil {
		return []models.Track{}, trackArtistsUpsertError
	}

	// Use type assertion to convert upsertedTracks to []models.Track
	tracks, assertError := assertTracks(upsertedTracks)
	if assertError != nil {
		return []models.Track{}, errors.New("could not convert upsertedTracks to []models.Track")
	}

	return tracks, nil
}
