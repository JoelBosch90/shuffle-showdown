package verification

import (
	models "api/database/models"
	"api/lib/wikipedia"
	"math"
	"time"

	"github.com/jinzhu/gorm"
)

func VerifyTrack(database *gorm.DB, trackId string) error {
	track := models.Track{}
	database.Model(&track).Preload("Artists").Where("id = ?", trackId).First(&track)

	artistNames := []string{}
	for _, artist := range track.Artists {
		artistNames = append(artistNames, artist.Name)
	}

	newReleaseYear, newReleaseYearError := wikipedia.GetTrackReleaseYear(track.Name, artistNames)
	if newReleaseYearError != nil {
		return newReleaseYearError
	}

	oldestReleaseYear := uint(math.Min(float64(newReleaseYear), float64(track.ReleaseYear)))
	if oldestReleaseYear != track.ReleaseYear {
		return database.Model(&track).Updates(models.Track{
			ReleaseYear: oldestReleaseYear,
			VerifiedAt:  time.Now(),
		}).Error
	}

	return database.Model(&track).Updates(models.Track{
		VerifiedAt: time.Now(),
	}).Error
}
