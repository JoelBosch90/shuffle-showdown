package game

import (
	"api/database"
	"api/database/models"
	"errors"

	uuid "github.com/satori/go.uuid"
)

func AwardInitialTracks(gameId uuid.UUID) error {
	database := database.Get()
	var game models.Game

	gameLoadError := database.Preload("Players").Where("id = ?", gameId).First(&game).Error
	if gameLoadError != nil {
		return errors.New("could not load game")
	}

	for _, player := range game.Players {
		track, selectTrackError := SelectNextTrack(gameId)
		if selectTrackError != nil {
			return selectTrackError
		}

		awardTrackError := AwardTrack(gameId, track, player)
		if awardTrackError != nil {
			return awardTrackError
		}
	}

	return nil
}
