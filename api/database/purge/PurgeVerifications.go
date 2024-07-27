package purge

import "github.com/jinzhu/gorm"

func PurgeVerifications(database *gorm.DB) error {
	purgeTrackVerificationError := database.Exec("DELETE FROM track_verifications WHERE NOT EXISTS (SELECT 1 FROM tracks WHERE track_verifications.track_id = tracks.id);").Error
	if purgeTrackVerificationError != nil {
		return purgeTrackVerificationError
	}

	return nil
}
