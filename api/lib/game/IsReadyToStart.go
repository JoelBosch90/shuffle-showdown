package game

import (
	database "api/database"
	"api/database/models"

	uuid "github.com/satori/go.uuid"
)

type CheckedTracksCount struct {
	TracksChecked uint `json:"tracks_checked"`
}

func IsReadyToStart(gameId uuid.UUID) bool {
	database := database.Get()

	var game models.Game
	loadGameError := database.Preload("Players").Where("id = ?", gameId).First(&game).Error
	if loadGameError != nil {
		return false
	}

	playerCount := uint(len(game.Players))
	checkedTracksCount := CheckedTracksCount{}
	countsError := database.Raw(`
		SELECT
			SUM(CASE WHEN check_completed_at IS NOT NULL THEN 1 ELSE 0 END) AS tracks_checked
		FROM tracks
		JOIN playlist_tracks ON tracks.id = playlist_tracks.track_id
		WHERE playlist_tracks.playlist_id = ?
	`, game.PlaylistId).Scan(&checkedTracksCount).Error
	if countsError != nil {
		return false
	}

	return checkedTracksCount.TracksChecked >= playerCount*game.SongsToWin
}
