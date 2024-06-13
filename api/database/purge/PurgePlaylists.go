package purge

import (
	"github.com/jinzhu/gorm"
)

func PurgePlaylists(database *gorm.DB) error {
	purgePlaylistError := database.Exec("DELETE FROM playlists WHERE NOT EXISTS (SELECT 1 FROM games WHERE playlists.id = games.playlist_id);").Error
	if purgePlaylistError != nil {
		return purgePlaylistError
	}

	purgePlaylistTrackError := database.Exec("DELETE FROM playlist_tracks WHERE NOT EXISTS (SELECT 1 FROM playlists WHERE playlist_tracks.playlist_id = playlists.id);").Error
	if purgePlaylistTrackError != nil {
		return purgePlaylistTrackError
	}

	purgeTrackError := database.Exec("DELETE FROM tracks WHERE NOT EXISTS (SELECT 1 FROM playlist_tracks WHERE tracks.id = playlist_tracks.track_id);").Error
	if purgeTrackError != nil {
		return purgeTrackError
	}

	purgeTrackArtistError := database.Exec("DELETE FROM track_artists WHERE NOT EXISTS (SELECT 1 FROM tracks WHERE track_artists.track_id = tracks.id);").Error
	if purgeTrackArtistError != nil {
		return purgeTrackArtistError
	}

	purgeArtistError := database.Exec("DELETE FROM artists WHERE NOT EXISTS (SELECT 1 FROM track_artists WHERE artists.id = track_artists.artist_id);").Error
	if purgeArtistError != nil {
		return purgeArtistError
	}

	return nil
}
