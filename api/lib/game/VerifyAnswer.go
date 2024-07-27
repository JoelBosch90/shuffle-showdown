package game

import (
	"api/database/models"
	"errors"

	uuid "github.com/satori/go.uuid"
)

func wonByPlayer(playerId uuid.UUID, wonTracks *[]models.WonTrack) []models.Track {
	wonByPlayer := []models.Track{}

	for _, wonTrack := range *wonTracks {
		if wonTrack.PlayerId == playerId {
			wonByPlayer = append(wonByPlayer, wonTrack.Track)
		}
	}

	return wonByPlayer
}

func getReleaseYears(tracks []models.Track) []uint {
	years := []uint{}

	for _, track := range tracks {
		years = append(years, track.ReleaseYear)
	}

	return years
}

func isInBetween(value uint, min uint, max uint) bool {
	return value > min && value < max
}

func hasBefore(values []uint, max uint) bool {
	for _, value := range values {
		if value < max {
			return true
		}
	}

	return false
}

func hasAfter(values []uint, min uint) bool {
	for _, value := range values {
		if value > min {
			return true
		}
	}

	return false
}

func hasInBetween(values []uint, min *int, max *int) bool {
	if min == nil && max == nil {
		return false
	}

	if min != nil && max == nil {
		minAsUint := uint(*min)
		return hasAfter(values, minAsUint)
	}

	if min == nil && max != nil {
		maxAsUint := uint(*max)
		return hasBefore(values, maxAsUint)
	}

	for _, value := range values {
		maxAsUint := uint(*max)
		minAsUint := uint(*min)
		if isInBetween(value, minAsUint, maxAsUint) {
			return true
		}
	}

	return false
}

func VerifyAnswer(afterReleaseYear *int, beforeReleaseYear *int, playerId uuid.UUID, rounds *[]models.Round, wonTracks *[]models.WonTrack) (bool, error) {
	if beforeReleaseYear == nil && afterReleaseYear == nil {
		return false, errors.New("answer must contain at least one field")
	}

	_, lastRound := LastRound(*rounds)
	if lastRound.PlayerId != playerId {
		return false, errors.New("not your turn")
	}

	wonTracksByPlayer := wonByPlayer(playerId, wonTracks)
	wonReleaseYears := getReleaseYears(wonTracksByPlayer)
	isCheating := hasInBetween(wonReleaseYears, afterReleaseYear, beforeReleaseYear)

	if isCheating {
		return false, errors.New("cheating is not allowed")
	}

	trackToWin := lastRound.Track
	correctBefore := beforeReleaseYear == nil || trackToWin.ReleaseYear <= uint(*beforeReleaseYear)
	correctAfter := afterReleaseYear == nil || trackToWin.ReleaseYear >= uint(*afterReleaseYear)

	return correctBefore && correctAfter, nil
}
