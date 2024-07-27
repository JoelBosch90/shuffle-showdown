package purge

import (
	"api/database/models"
	"time"

	"github.com/jinzhu/gorm"
)

func PurgeAccessTokens(database *gorm.DB) error {
	return database.Delete(&models.SpotifyAccessToken{}, "expires_at < ?", time.Now()).Error
}
