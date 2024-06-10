package spotify

import (
	databaseHelpers "api/database/helpers"
	"api/database/models"
	spotifyModels "api/lib/spotify/models"

	"github.com/jinzhu/gorm"
)

func CreatePlaylist(info spotifyModels.Playlist, countryCode string, database *gorm.DB) (models.Playlist, error) {
	lastSongAdded, tracks, tracksError := CreateTracks(database, info.Tracks.Items)
	if tracksError != nil {
		return models.Playlist{}, tracksError
	}

	upsertedPlaylist, upsertPlaylistError := databaseHelpers.Upsert(database, []interface{}{&models.Playlist{
		Id:            info.Id,
		Name:          info.Name,
		CountryCode:   countryCode,
		LastSongAdded: lastSongAdded,
		TracksTotal:   uint(info.Tracks.Total),
		Tracks:        tracks,
	}})
	if upsertPlaylistError != nil || len(upsertedPlaylist) == 0 {
		return models.Playlist{}, upsertPlaylistError
	}

	playlistTrackLinks := []interface{}{}
	for _, track := range tracks {
		playlistTrackLinks = append(playlistTrackLinks, &models.PlaylistTrack{
			PlaylistId: info.Id,
			TrackId:    track.Id,
		})
	}

	_, upsertPlayListTrackLinksError := databaseHelpers.Upsert(database, playlistTrackLinks)
	if upsertPlayListTrackLinksError != nil {
		return models.Playlist{}, upsertPlayListTrackLinksError
	}

	var playlist models.Playlist
	getPlaylistError := database.Where("id = ?", info.Id).First(&playlist).Error
	if getPlaylistError != nil {
		return models.Playlist{}, getPlaylistError
	}

	return playlist, nil
}
