package game

import (
	database "api/database"
	"api/database/models"
	"api/lib/helpers"
	"errors"
	"math/rand"

	uuid "github.com/satori/go.uuid"
)

func pickRandomTrack(tracks []models.Track) models.Track {
	return tracks[rand.Intn(len(tracks))]
}

func SelectNextTrack(gameId uuid.UUID) (models.Track, error) {
	database := database.Get()
	var game models.Game

	loadGameError := database.Preload("WonTracks").Preload("Playlist").Preload("Playlist.Tracks", "check_completed_at IS NOT NULL").Where("id = ?", gameId).First(&game).Error
	if loadGameError != nil {
		return models.Track{}, errors.New("could not load game")
	}

	wonTrackIds := []string{}
	for _, wonTrack := range game.WonTracks {
		wonTrackIds = append(wonTrackIds, wonTrack.TrackId)
	}

	availableTracks := []models.Track{}
	for _, track := range game.Playlist.Tracks {
		if !helpers.IncludesString(wonTrackIds, track.Id) {
			availableTracks = append(availableTracks, track)
		}
	}

	if len(availableTracks) <= 0 {
		return models.Track{}, errors.New("no tracks left")
	}

	return pickRandomTrack(availableTracks), nil
}
