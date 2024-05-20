package game

import (
	"api/database"
	models "api/database/models"

	uuid "github.com/satori/go.uuid"
)

func AwardTrack(gameId uuid.UUID, track models.Track, player models.Player) error {
	database := database.Get()
	result := database.Where("game_id = ? AND track_id = ? AND player_id = ?", gameId, track.Id, player.Id).Limit(1).Find(&models.WonTrack{})
	// Never award the same track to the same player twice
	if result.RowsAffected != 0 {
		return nil
	}
	if result.Error != nil && result.Error.Error() != "record not found" {
		return result.Error
	}

	return database.Create(&models.WonTrack{
		Id:       uuid.NewV4(),
		GameId:   gameId,
		TrackId:  track.Id,
		PlayerId: player.Id,
	}).Error
}
