package spotify

import (
	databaseHelpers "api/database"
	"api/database/models"
	spotifyModels "api/lib/spotify/models"
	"api/lib/verification"
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

func findTrackById(tracks []models.Track, id string) *models.Track {
	for _, track := range tracks {
		if track.Id == id {
			return &track
		}
	}

	return nil
}

func getExistingTracks(database *gorm.DB, items []spotifyModels.Item) []models.Track {
	var existingTracks []models.Track
	var trackIds []string

	for _, item := range items {
		trackIds = append(trackIds, item.Track.Id)
	}

	database.Where("id IN (?)", trackIds).Find(&existingTracks)

	return existingTracks
}

func getOldestReleaseYear(newTrack *models.Track, existingTrack *models.Track) uint {
	if existingTrack == nil {
		return newTrack.ReleaseYear
	}

	if newTrack.ReleaseYear > existingTrack.ReleaseYear {
		return existingTrack.ReleaseYear
	}

	return newTrack.ReleaseYear
}

func getOldestReleaseMonth(newTrack *models.Track, existingTrack *models.Track) uint {
	if existingTrack == nil {
		return newTrack.ReleaseMonth
	}

	if newTrack.ReleaseYear < existingTrack.ReleaseYear {
		return newTrack.ReleaseMonth
	}

	if newTrack.ReleaseYear == existingTrack.ReleaseYear && newTrack.ReleaseMonth > existingTrack.ReleaseMonth {
		return existingTrack.ReleaseMonth
	}

	return newTrack.ReleaseMonth
}

func getOldestReleaseDay(newTrack *models.Track, existingTrack *models.Track) uint {
	if existingTrack == nil {
		return newTrack.ReleaseDay
	}

	if newTrack.ReleaseYear < existingTrack.ReleaseYear {
		return newTrack.ReleaseDay
	}

	if newTrack.ReleaseYear == existingTrack.ReleaseYear && newTrack.ReleaseMonth < existingTrack.ReleaseMonth {
		return newTrack.ReleaseDay
	}

	if newTrack.ReleaseYear == existingTrack.ReleaseYear && newTrack.ReleaseMonth == existingTrack.ReleaseMonth && newTrack.ReleaseDay > existingTrack.ReleaseDay {
		return existingTrack.ReleaseDay
	}

	return newTrack.ReleaseDay
}

func constructTracks(database *gorm.DB, items []spotifyModels.Item, createdArtists []models.Artist) ([]interface{}, []interface{}, error) {
	var tracksToCreate []interface{}
	var trackArtistsToCreate []interface{}
	existingTracks := getExistingTracks(database, items)

	for _, item := range items {
		trackToCreate := item.Track

		releaseYear, releaseMonth, releaseDay := ConvertReleaseDateToIntegers(item.Album.ReleaseDate)

		if releaseYear == 0 || trackToCreate.PreviewUrl == "" {
			continue
		}

		artistIds := []string{}
		for _, artist := range trackToCreate.Artists {
			artistIds = append(artistIds, artist.Id)
		}

		artists := []models.Artist{}
		for _, artistId := range artistIds {
			artist := findArtistById(createdArtists, artistId)
			if artist != nil {
				artists = append(artists, *artist)
			}

			trackArtistsToCreate = append(trackArtistsToCreate, &models.TrackArtist{
				TrackId:  trackToCreate.Id,
				ArtistId: artistId,
			})
		}

		newDates := &models.Track{
			ReleaseYear:  releaseYear,
			ReleaseMonth: releaseMonth,
			ReleaseDay:   releaseDay,
		}

		existingTrack := findTrackById(existingTracks, trackToCreate.Id)
		trackUpdate := &models.Track{
			Id:           trackToCreate.Id,
			Name:         trackToCreate.Name,
			ReleaseYear:  getOldestReleaseYear(newDates, existingTrack),
			ReleaseMonth: getOldestReleaseMonth(newDates, existingTrack),
			ReleaseDay:   getOldestReleaseDay(newDates, existingTrack),
			Artists:      artists,
			PreviewUrl:   trackToCreate.PreviewUrl,
			IsPlayable:   trackToCreate.IsPlayable,
		}

		tracksToCreate = append(tracksToCreate, trackUpdate)
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
	var tracks []models.Track
	var assertError error

	createError := database.Transaction(func(transaction *gorm.DB) error {
		createdArtists, createdArtistError := CreateArtists(database, items)
		if createdArtistError != nil || len(createdArtists) == 0 {
			return createdArtistError
		}

		tracksToCreate, trackArtistsToCreate, constructError := constructTracks(database, items, createdArtists)
		if constructError != nil || len(tracksToCreate) == 0 {
			return constructError
		}

		upsertedTracks, upsertError := databaseHelpers.Upsert(database, tracksToCreate)
		if upsertError != nil || len(upsertedTracks) == 0 {
			return upsertError
		}

		_, trackArtistsUpsertError := databaseHelpers.Upsert(database, trackArtistsToCreate)
		if trackArtistsUpsertError != nil {
			return trackArtistsUpsertError
		}

		// Use type assertion to convert upsertedTracks to []models.Track
		tracks, assertError = assertTracks(upsertedTracks)
		if assertError != nil {
			return errors.New("could not convert upsertedTracks to []models.Track")
		}

		for _, track := range tracks {
			verificationError := verification.CreateTrackVerification(database, track.Id)
			if verificationError != nil {
				return verificationError
			}
		}

		return nil
	})
	if createError != nil {
		return []models.Track{}, createError
	}

	return tracks, nil
}
