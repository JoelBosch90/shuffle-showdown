package game

import (
	"api/database/models"
	databaseHelpers "api/lib/database_helpers"
	spotifyModels "api/lib/spotify_models"
	"errors"

	"github.com/jinzhu/gorm"
)

func CreatePlaylist(info spotifyModels.Playlist, database *gorm.DB) (models.Playlist, error) {
	lastSongAdded, tracks, tracksError := CreateTracks(database, info.Tracks.Items)
	if tracksError != nil {
		return models.Playlist{}, tracksError
	}

	upsertedPlaylist, upsertPlaylistError := databaseHelpers.Upsert(database, []interface{}{&models.Playlist{
		Id:            info.Id,
		Name:          info.Name,
		LastSongAdded: lastSongAdded,
		TracksTotal:   uint(info.Tracks.Total),
		Tracks:        tracks,
	}})
	if upsertPlaylistError != nil || len(upsertedPlaylist) == 0 {
		return models.Playlist{}, upsertPlaylistError
	}

	// Create links between the playlist and the tracks on the playlist.
	playlistTrackLinks := []interface{}{}
	for _, track := range tracks {
		playlistTrackLinks = append(playlistTrackLinks, &models.PlaylistTrack{
			PlaylistId: info.Id,
			TrackId:    track.Id,
		})
	}

	// Upsert the playlist track links.
	_, upsertPlayListTrackLinksError := databaseHelpers.Upsert(database, playlistTrackLinks)
	if upsertPlayListTrackLinksError != nil {
		return models.Playlist{}, upsertPlayListTrackLinksError
	}

	// Use type assertion to convert upsertedPlaylist[0] to models.Playlist
	playlist, ok := upsertedPlaylist[0].(models.Playlist)
	if !ok {
		return models.Playlist{}, errors.New("could not convert upsertedPlaylist[0] to models.Playlist")
	}

	return playlist, nil
}
