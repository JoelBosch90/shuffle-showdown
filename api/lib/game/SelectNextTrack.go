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

	loadGameError := database.Preload("Playlist.Tracks").Preload("Rounds").Where("id = ?", gameId).First(&game).Error
	if loadGameError != nil {
		return models.Track{}, errors.New("could not load game")
	}
	if len(game.Playlist.Tracks) < int(game.SongsToWin) {
		return models.Track{}, errors.New("too few tracks")
	}

	playedTrackIds := []string{}
	for _, round := range game.Rounds {
		playedTrackIds = append(playedTrackIds, round.TrackId)
	}

	availableTracks := []models.Track{}
	for _, track := range game.Playlist.Tracks {
		if !includesString(playedTrackIds, track.Id) {
			availableTracks = append(availableTracks, track)
		}
	}

	return availableTracks[rand.Intn(len(availableTracks))], nil
}
