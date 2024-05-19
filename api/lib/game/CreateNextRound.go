package game

import (
	database "api/database"
	"api/database/models"
	"errors"

	uuid "github.com/satori/go.uuid"
)

func findNextRoundNumber(rounds []models.Round) uint {
	if len(rounds) == 0 {
		return 1
	}

	return rounds[len(rounds)-1].Number + 1
}

func CreateNextRound(gameId uuid.UUID) error {
	database := database.Get()
	var game models.Game

	loadGameError := database.Preload("Playlist.Tracks").Preload("Rounds").Where("id = ?", gameId).First(&game).Error
	if loadGameError != nil {
		return errors.New("could not load game")
	}

	nextPlayerId, findNextPlayerIdError := FindNextPlayerId(game.Id)
	if findNextPlayerIdError != nil {
		return findNextPlayerIdError
	}

	nextTrack, findNextTrackError := SelectNextTrack(game.Id)
	if findNextTrackError != nil {
		return findNextTrackError
	}

	return database.Create(&models.Round{
		Id:       uuid.NewV4(),
		Number:   findNextRoundNumber(game.Rounds),
		GameId:   gameId,
		PlayerId: nextPlayerId,
		TrackId:  nextTrack.Id,
	}).Error
}
