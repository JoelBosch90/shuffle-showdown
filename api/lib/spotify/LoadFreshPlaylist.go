package spotify

import (
	"api/database"
	"api/database/models"
	"errors"
	"time"
)

var PLAYLIST_CACHE_TIME = time.Hour * 24

func getRecentPlaylistFromDatabase(playlistId string, countryCode string) (models.Playlist, error) {
	database := database.Get()

	playlist := models.Playlist{}
	playlistError := database.Model(&playlist).Where("id = ?", playlistId).First(&playlist).Error
	if playlistError != nil {
		return models.Playlist{}, playlistError
	}

	if playlist.CountryCode != countryCode {
		return models.Playlist{}, errors.New("playlist not found")
	}

	if playlist.UpdatedAt.Before(time.Now().Add(-PLAYLIST_CACHE_TIME)) {
		return models.Playlist{}, errors.New("playlist too old")
	}

	return playlist, nil
}

func LoadFreshPlaylist(playlistId string, countryCode string, sendLoadingUpdate func(tracksLoaded int, tracksTotal int)) error {
	database := database.Get()

	_, playlistError := getRecentPlaylistFromDatabase(playlistId, countryCode)
	if playlistError == nil {
		return nil
	}

	info, infoError := RequestPlaylistInfo(playlistId, countryCode, sendLoadingUpdate)
	if infoError != nil {
		return infoError
	}

	_, newPlaylistError := CreatePlaylist(info, countryCode, database)
	if newPlaylistError != nil {
		return newPlaylistError
	}

	return nil
}
