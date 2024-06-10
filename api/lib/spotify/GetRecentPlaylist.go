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

	if playlist.CreatedAt.Before(time.Now().Add(-PLAYLIST_CACHE_TIME)) {
		return models.Playlist{}, errors.New("playlist too old")
	}

	return playlist, nil
}

func GetRecentPlaylist(playlistId string, countryCode string) (models.Playlist, error) {
	database := database.Get()

	playlist, playlistError := getRecentPlaylistFromDatabase(playlistId, countryCode)
	if playlistError == nil {
		return playlist, nil
	}

	info, infoError := RequestPlaylistInfo(playlistId, countryCode)
	if infoError != nil {
		return models.Playlist{}, infoError
	}

	newPlaylist, newPlaylistError := CreatePlaylist(info, countryCode, database)
	if newPlaylistError != nil {
		return models.Playlist{}, newPlaylistError
	}

	return newPlaylist, nil
}
