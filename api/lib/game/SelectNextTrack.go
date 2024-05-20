package game

import (
	database "api/database"
	"api/database/models"
	"errors"
	"math/rand"
	"slices"

	uuid "github.com/satori/go.uuid"
)

func includesString(haystack []string, needle string) bool {
	for _, text := range haystack {
		if text == needle {
			return true
		}
	}
	return false
}

func withoutDuplicates(withDuplicates []string) []string {
	unique := withDuplicates
	slices.Sort(unique)

	return slices.Compact(unique)
}

func SelectNextTrack(gameId uuid.UUID) (models.Track, error) {
	database := database.Get()
	var game models.Game

	loadGameError := database.Preload("Rounds").Preload("WonTracks").Preload("Playlist").Preload("Playlist.Tracks").Where("id = ?", gameId).First(&game).Error
	if loadGameError != nil {
		return models.Track{}, errors.New("could not load game")
	}

	usedTrackIds := []string{}
	for _, round := range game.Rounds {
		usedTrackIds = append(usedTrackIds, round.TrackId)
	}
	for _, wonTrack := range game.WonTracks {
		usedTrackIds = append(usedTrackIds, wonTrack.TrackId)
	}

	usedTrackIds = withoutDuplicates(usedTrackIds)
	availableTracks := []models.Track{}
	for _, track := range game.Playlist.Tracks {
		if !includesString(usedTrackIds, track.Id) {
			availableTracks = append(availableTracks, track)
		}
	}

	if len(availableTracks) <= 0 {
		return models.Track{}, errors.New("no tracks left")
	}

	return availableTracks[rand.Intn(len(availableTracks))], nil
}
