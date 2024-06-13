package purge

import (
	"api/database/models"
	"time"

	"github.com/jinzhu/gorm"
)

func PurgeAccessTokens(database *gorm.DB) error {
	return database.Delete(&models.AccessToken{}, "expires_at < ?", time.Now()).Error
}
