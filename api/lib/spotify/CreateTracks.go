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

func constructTracks(items []spotifyModels.Item, createdArtists []models.Artist) ([]interface{}, []interface{}, error) {
	var tracksToCreate []interface{}
	var trackArtistsToCreate []interface{}

	for _, item := range items {
		track := item.Track

		releaseYear, releaseMonth, releaseDay := ConvertReleaseDateToIntegers(track.Album.ReleaseDate)

		if releaseYear == 0 || track.PreviewUrl == "" {
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
			Id:           track.Id,
			Name:         track.Name,
			ReleaseYear:  releaseYear,
			ReleaseMonth: releaseMonth,
			ReleaseDay:   releaseDay,
			Artists:      artists,
			PreviewUrl:   track.PreviewUrl,
			IsPlayable:   track.IsPlayable,
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

func CreateTracks(database *gorm.DB, items []spotifyModels.Item) (string, []models.Track, error) {
	var artistsToCreate []spotifyModels.Artist
	var lastSongAdded string = ""

	for _, item := range items {
		lastSongAdded = item.AddedAt
		artistsToCreate = append(artistsToCreate, item.Track.Artists...)
	}

	createdArtists, createdArtistError := CreateArtists(database, artistsToCreate)
	if createdArtistError != nil || len(createdArtists) == 0 {
		return lastSongAdded, []models.Track{}, createdArtistError
	}

	tracksToCreate, trackArtistsToCreate, constructError := constructTracks(items, createdArtists)
	if constructError != nil || len(tracksToCreate) == 0 {
		return lastSongAdded, []models.Track{}, constructError
	}

	upsertedTracks, upsertError := databaseHelpers.Upsert(database, tracksToCreate)
	if upsertError != nil || len(upsertedTracks) == 0 {
		return lastSongAdded, []models.Track{}, upsertError
	}

	_, trackArtistsUpsertError := databaseHelpers.Upsert(database, trackArtistsToCreate)
	if trackArtistsUpsertError != nil {
		return lastSongAdded, []models.Track{}, trackArtistsUpsertError
	}

	// Use type assertion to convert upsertedTracks to []models.Track
	tracks, assertError := assertTracks(upsertedTracks)
	if assertError != nil {
		return lastSongAdded, []models.Track{}, errors.New("could not convert upsertedTracks to []models.Track")
	}

	return lastSongAdded, tracks, nil
}
