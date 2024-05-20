package game

import (
	database "api/database"
	"api/database/models"
	"errors"
	"math/rand"

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

func SelectNextTrack(gameId uuid.UUID) (models.Track, error) {
	database := database.Get()
	var game models.Game

	loadGameError := database.Preload("Playlist.Tracks").Preload("Rounds").Preload("WonTracks").Where("id = ?", gameId).First(&game).Error
	if loadGameError != nil {
		return models.Track{}, errors.New("could not load game")
	}
	if len(game.Playlist.Tracks) < int(game.SongsToWin) {
		return models.Track{}, errors.New("too few tracks")
	}

	usedTrackIds := []string{}
	for _, round := range game.Rounds {
		usedTrackIds = append(usedTrackIds, round.TrackId)
	}
	for _, wonTrack := range game.WonTracks {
		usedTrackIds = append(usedTrackIds, wonTrack.TrackId)
	}

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
