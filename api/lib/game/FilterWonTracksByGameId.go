package game

import (
	"api/database/models"

	uuid "github.com/satori/go.uuid"
)

func FilterWonTracksByGameId(gameId uuid.UUID, wonTracks []models.WonTrack) []models.WonTrack {
	filtered := []models.WonTrack{}
	for _, wonTrack := range wonTracks {
		if wonTrack.GameId == gameId {
			filtered = append(filtered, wonTrack)
		}
	}

	return filtered
}
