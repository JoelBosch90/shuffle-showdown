package verification

import (
	"api/database"
	models "api/database/models"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

var MAX_PROCESSING_TIME = time.Minute * 5
var RETRY_INTERVAL = time.Millisecond * 200

func isVerificationInProgress(database *gorm.DB) bool {
	maxProcessingTimeAgo := time.Now().Add(-MAX_PROCESSING_TIME)
	return !database.Model(&models.Track{}).Where("(check_started_at IS NOT NULL OR check_started_at < ?) AND check_completed_at IS NULL", maxProcessingTimeAgo).First(&models.Track{}).RecordNotFound()
}

func startCheck(database *gorm.DB, trackId string) error {
	return database.Model(&models.Track{}).Where("id = ?", trackId).Update("check_started_at", time.Now()).Error
}

func completeCheck(database *gorm.DB, trackId string) error {
	return database.Model(&models.Track{}).Where("id = ?", trackId).Update("check_completed_at", time.Now()).Error
}

func fetchOldestUnverifiedTrack(database *gorm.DB) (models.Track, error) {
	oldestUnverifiedTrack := models.Track{}
	fetchError := database.Model(&models.Track{}).Where("check_started_at IS NULL AND check_completed_at IS NULL").Order("created_at").First(&oldestUnverifiedTrack).Error

	if oldestUnverifiedTrack.Id == "" {
		return models.Track{}, errors.New("no unverified tracks")
	}

	return oldestUnverifiedTrack, fetchError
}

func VerifyTracks() error {
	database := database.Get()
	verifiedTrackId := ""

	for {
		transactionError := database.Transaction(func(transaction *gorm.DB) error {
			if verifiedTrackId != "" {
				completeError := completeCheck(transaction, verifiedTrackId)
				if completeError != nil {
					return completeError
				}
			}

			if isVerificationInProgress(transaction) {
				return errors.New("verification in progress")
			}

			oldestUnverifiedTrack, fetchError := fetchOldestUnverifiedTrack(transaction)
			if fetchError != nil {
				return fetchError
			}

			verifiedTrackId = oldestUnverifiedTrack.Id
			startError := startCheck(transaction, verifiedTrackId)
			if startError != nil {
				return startError
			}

			return nil
		})

		if transactionError != nil {
			if transactionError.Error() == "database is locked" {
				time.Sleep(RETRY_INTERVAL)
				continue
			}

			if verifiedTrackId != "" {
				return completeCheck(database, verifiedTrackId)
			}

			return nil
		}

		VerifyTrack(database, verifiedTrackId)
	}
}
