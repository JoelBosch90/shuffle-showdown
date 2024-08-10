package verification

import (
	"api/database"
	models "api/database/models"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

var MAX_PROCESSING_TIME = time.Minute * 5
var RETRY_INTERVAL = time.Millisecond * 25

func isVerificationInProgress(database *gorm.DB, now time.Time) bool {
	maxProcessingTimeAgo := now.Add(-MAX_PROCESSING_TIME)
	trackInProgress := models.Track{}
	databaseError := database.Model(&models.Track{}).Where("(check_started_at IS NOT NULL OR check_started_at >= ?) AND check_completed_at IS NULL", maxProcessingTimeAgo).First(&trackInProgress).Error
	if errors.Is(databaseError, gorm.ErrRecordNotFound) {
		return false
	}
	return databaseError == nil
}

func startCheck(database *gorm.DB, trackId string, now time.Time) error {
	return database.Model(&models.Track{}).Where("id = ?", trackId).Update("check_started_at", now).Error
}

func completeCheck(database *gorm.DB, trackId string) error {
	return database.Model(&models.Track{}).Where("id = ?", trackId).Update("check_completed_at", time.Now()).Error
}

func fetchOldestUnverifiedTrack(database *gorm.DB, now time.Time) (models.Track, error) {
	maxProcessingTimeAgo := now.Add(-MAX_PROCESSING_TIME)
	oldestUnverifiedTrack := models.Track{}
	fetchError := database.Model(&models.Track{}).Where("(check_started_at IS NULL OR check_started_at < ?) AND check_completed_at IS NULL", maxProcessingTimeAgo).Order("created_at").First(&oldestUnverifiedTrack).Error

	if oldestUnverifiedTrack.Id == "" {
		return models.Track{}, errors.New("no unverified tracks")
	}

	return oldestUnverifiedTrack, fetchError
}

func VerifyTracks(sendCheckingUpdate func()) error {
	database := database.Get()
	unverifiedTrackId := ""
	verifiedTrackId := ""

	for {
		transactionError := database.Transaction(func(transaction *gorm.DB) error {
			if verifiedTrackId != "" {
				completeError := completeCheck(transaction, verifiedTrackId)
				if completeError != nil {
					return completeError
				}
				sendCheckingUpdate()
			}

			checkStartTimeStamp := time.Now()
			inProgress := isVerificationInProgress(transaction, checkStartTimeStamp)
			if inProgress {
				return errors.New("verification in progress")
			}

			oldestUnverifiedTrack, fetchError := fetchOldestUnverifiedTrack(transaction, checkStartTimeStamp)
			if fetchError != nil {
				return fetchError
			}

			unverifiedTrackId = oldestUnverifiedTrack.Id
			startError := startCheck(transaction, unverifiedTrackId, checkStartTimeStamp)
			if startError != nil {
				return startError
			}

			return nil
		})

		if transactionError != nil {
			if transactionError.Error() == "database is locked" {
				unverifiedTrackId = ""
				time.Sleep(RETRY_INTERVAL)
				continue
			}

			if verifiedTrackId != "" {
				completeError := completeCheck(database, verifiedTrackId)
				if completeError == nil {
					sendCheckingUpdate()
				}
				return completeError
			}

			return nil
		}

		VerifyTrack(database, unverifiedTrackId)
		verifiedTrackId = unverifiedTrackId
	}
}
