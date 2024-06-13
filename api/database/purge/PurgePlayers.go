package purge

import (
	"time"

	"github.com/jinzhu/gorm"
)

func PurgePlayers(database *gorm.DB, maxAgeInDays int) error {
	thresholdDate := time.Now().AddDate(0, 0, -maxAgeInDays)

	findPlayerIdsError := database.Exec("DELETE FROM players WHERE \"created_at\" < ? AND NOT EXISTS (SELECT 1 FROM game_players WHERE players.id = game_players.player_id);", thresholdDate).Error
	if findPlayerIdsError != nil {
		return findPlayerIdsError
	}

	return nil
}
