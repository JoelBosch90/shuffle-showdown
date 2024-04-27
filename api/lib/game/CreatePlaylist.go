package game

import (
	"api/database/models"
	databaseHelpers "api/lib/database_helpers"
	spotifyModels "api/lib/spotify_models"

	"github.com/jinzhu/gorm"
)

func CreatePlaylist(info spotifyModels.Playlist, database *gorm.DB) error {
	var lastSongAdded string = ""

	// TODO: Add tracks
	// TODO: Add lastSongAdded

	upsertPlaylistError := databaseHelpers.Upsert(database, []interface{}{&models.Playlist{
		Id:            info.Id,
		Name:          info.Name,
		LastSongAdded: lastSongAdded,
		TracksTotal:   uint(info.Tracks.Total),
	}})
	if upsertPlaylistError != nil {
		return upsertPlaylistError
	}

	return nil
}
