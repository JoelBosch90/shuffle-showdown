package purge

import "github.com/jinzhu/gorm"

func PurgeAll(database *gorm.DB, maxAgeInDays int) error {
	purgeAccessTokensError := PurgeAccessTokens(database)
	if purgeAccessTokensError != nil {
		return purgeAccessTokensError
	}

	purgeGamesError := PurgeGames(database, maxAgeInDays)
	if purgeGamesError != nil {
		return purgeGamesError
	}

	purgePlayersError := PurgePlayers(database, maxAgeInDays)
	if purgePlayersError != nil {
		return purgePlayersError
	}

	purgePlaylistsError := PurgePlaylists(database)
	if purgePlaylistsError != nil {
		return purgePlaylistsError
	}

	return nil
}
