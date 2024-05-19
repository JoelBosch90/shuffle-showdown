package game

import (
	"api/database"
	models "api/database/models"

	uuid "github.com/satori/go.uuid"
)

func AwardTrack(gameId uuid.UUID, track models.Track, player models.Player) error {
	database := database.Get()
	return database.Create(&models.WonTrack{
		Id:       uuid.NewV4(),
		GameId:   gameId,
		TrackId:  track.Id,
		PlayerId: player.Id,
	}).Error
}
