package game

import (
	"api/database/models"
	databaseHelpers "api/lib/database_helpers"
	spotifyModels "api/lib/spotify_models"
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

	// Loop through the items to construct the tracks.
	for _, item := range items {
		track := item.Track

		// Convert the release date to integers.
		releaseYear, releaseMonth, releaseDay := ConvertReleaseDateToIntegers(track.Album.ReleaseDate)

		// Find the ids artists to attach to the track.
		artistIds := []string{}
		for _, artist := range track.Artists {
			artistIds = append(artistIds, artist.Id)
		}

		// Find the corresponding artists to attach to the track.
		artists := []models.Artist{}
		for _, artistId := range artistIds {
			artist := findArtistById(createdArtists, artistId)
			if artist != nil {
				artists = append(artists, *artist)
			}

			// Also construct the links for the join table.
			trackArtistsToCreate = append(trackArtistsToCreate, &models.TrackArtist{
				TrackId:  track.Id,
				ArtistId: artistId,
			})
		}

		// Append the track to the tracksToCreate slice.
		tracksToCreate = append(tracksToCreate, &models.Track{
			Id:           track.Id,
			Name:         track.Name,
			ReleaseYear:  releaseYear,
			ReleaseMonth: releaseMonth,
			ReleaseDay:   releaseDay,
			Artists:      artists,
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

	// Loop through once to get the last time a song was added and all the artists to create.
	for _, item := range items {
		lastSongAdded = item.AddedAt
		artistsToCreate = append(artistsToCreate, item.Track.Artists...)
	}

	// Create the artists
	createdArtists, createdArtistError := CreateArtists(database, artistsToCreate)
	if createdArtistError != nil || len(createdArtists) == 0 {
		return lastSongAdded, []models.Track{}, createdArtistError
	}

	// Construct the tracks
	tracksToCreate, trackArtistsToCreate, constructError := constructTracks(items, createdArtists)
	if constructError != nil || len(tracksToCreate) == 0 {
		return lastSongAdded, []models.Track{}, constructError
	}

	// Upsert the tracks
	upsertedTracks, upsertError := databaseHelpers.Upsert(database, tracksToCreate)
	if upsertError != nil || len(upsertedTracks) == 0 {
		return lastSongAdded, []models.Track{}, upsertError
	}

	// Upsert the track artists
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
