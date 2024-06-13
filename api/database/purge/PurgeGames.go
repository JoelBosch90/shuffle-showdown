package purge

import (
	"api/database/models"
	"time"

	"github.com/jinzhu/gorm"
)

func PurgeGames(database *gorm.DB, maxAgeInDays int) error {
	thresholdDate := time.Now().AddDate(0, 0, -maxAgeInDays)
	purgeOldGamesError := database.Where("created_at < ?", thresholdDate).Delete(&models.Game{}).Error
	if purgeOldGamesError != nil {
		return purgeOldGamesError
	}

	purgeWonTracksError := database.Exec("DELETE FROM won_tracks WHERE NOT EXISTS (SELECT 1 FROM games WHERE won_tracks.game_id = games.id);").Error
	if purgeWonTracksError != nil {
		return purgeWonTracksError
	}

	purgeGamePlayersError := database.Exec("DELETE FROM game_players WHERE NOT EXISTS (SELECT 1 FROM games WHERE game_players.game_id = games.id);").Error
	if purgeGamePlayersError != nil {
		return purgeGamePlayersError
	}

	purgeRoundsError := database.Exec("DELETE FROM rounds WHERE NOT EXISTS (SELECT 1 FROM games WHERE rounds.game_id = games.id);").Error
	if purgeRoundsError != nil {
		return purgeRoundsError
	}

	purgeBannedPlayersError := database.Exec("DELETE FROM banned_players WHERE NOT EXISTS (SELECT 1 FROM games WHERE banned_players.game_id = games.id);").Error
	if purgeBannedPlayersError != nil {
		return purgeBannedPlayersError
	}

	return nil
}
