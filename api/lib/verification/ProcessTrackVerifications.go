package verification

import (
	"api/database"
	models "api/database/models"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

var MAX_PROCESSING_TIME = time.Hour

func hasRunningVerifications(database *gorm.DB) bool {
	thresholdTimestamp := time.Now().Add(-1 * MAX_PROCESSING_TIME)
	runningVerifications := []models.TrackVerification{}
	runningVerificationsError := database.Model(&runningVerifications).Where("started_at > ? AND completed_at IS NOT NULL", thresholdTimestamp).Error

	if runningVerificationsError != nil && runningVerificationsError.Error() == "record not found" {
		return false
	}

	if len(runningVerifications) == 0 {
		return false
	}

	return true
}

func getVerificationToProcess(database *gorm.DB) (models.TrackVerification, error) {
	verificationToProcess := models.TrackVerification{}
	zeroTimestamp := time.Time{}

	fetchError := database.Model(&verificationToProcess).Where("started_at = ? AND completed_at = ?", zeroTimestamp, zeroTimestamp).Order("created_at").First(&verificationToProcess).Error
	if fetchError != nil {
		return models.TrackVerification{}, fetchError
	}

	return verificationToProcess, nil
}

func processVerification(database *gorm.DB) (bool, error) {
	verificationToProcess, getError := getVerificationToProcess(database)
	if getError != nil || verificationToProcess.TrackId == "" {
		return true, getError
	}

	verificationToProcess.StartedAt = time.Now()
	updateError := database.Save(&verificationToProcess).Error
	if updateError != nil {
		return false, updateError
	}

	verificationError := VerifyTrack(database, verificationToProcess.TrackId)
	if verificationError != nil {
		return false, verificationError
	}

	verificationToProcess.CompletedAt = time.Now()
	updateError = database.Save(&verificationToProcess).Error
	if updateError != nil {
		return false, updateError
	}

	if time.Since(verificationToProcess.StartedAt) > MAX_PROCESSING_TIME {
		return false, errors.New("aborting because verification exceeded maximum processing time")
	}

	return false, nil
}

func ProcessTrackVerifications() error {
	database := database.Get()

	if hasRunningVerifications(database) {
		return nil
	}

	for finished, processingError := processVerification(database); processingError != nil; finished, processingError = processVerification(database) {
		if finished {
			return nil
		}
	}

	return nil
}
