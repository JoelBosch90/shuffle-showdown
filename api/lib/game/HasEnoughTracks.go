package game

import (
	database "api/database"
	"api/database/models"

	uuid "github.com/satori/go.uuid"
)

type TracksCount struct {
	Tracks uint `json:"tracks"`
}

func HasEnoughTracks(gameId uuid.UUID) bool {
	database := database.Get()

	var game models.Game
	loadGameError := database.Preload("Players").Where("id = ?", gameId).First(&game).Error
	if loadGameError != nil {
		return false
	}

	playerCount := uint(len(game.Players))
	tracksCount := TracksCount{}
	countsError := database.Raw(`
		SELECT
			COUNT(*) AS tracks
		FROM tracks
		JOIN playlist_tracks ON tracks.id = playlist_tracks.track_id
		WHERE playlist_tracks.playlist_id = ?
	`, game.PlaylistId).Scan(&tracksCount).Error
	if countsError != nil {
		return false
	}

	return tracksCount.Tracks >= playerCount*game.SongsToWin
}
