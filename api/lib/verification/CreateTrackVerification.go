package verification

import (
	models "api/database/models"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func CreateTrackVerification(database *gorm.DB, trackId string) error {
	existingVerification := models.TrackVerification{}
	existingVerificationError := database.Model(&existingVerification).Where("track_id = ?", trackId).First(&existingVerification).Error
	if existingVerificationError.Error() != "record not found" {
		return existingVerificationError
	}

	if !existingVerification.CompletedAt.IsZero() {
		return nil
	}

	createVerificationError := database.Create(&models.TrackVerification{
		Id:      uuid.NewV4(),
		TrackId: trackId,
	}).Error
	if createVerificationError != nil {
		return createVerificationError
	}

	go ProcessTrackVerifications()

	return nil
}
