package verification

import (
	models "api/database/models"
	"api/lib/wikipedia"
	"log"
	"math"

	"github.com/jinzhu/gorm"
)

func VerifyTrack(database *gorm.DB, trackId string) error {
	track := models.Track{}
	fetchError := database.Model(&track).Preload("Artists").Where("id = ?", trackId).First(&track).Error
	if fetchError != nil {
		return fetchError
	}

	artistNames := []string{}
	for _, artist := range track.Artists {
		artistNames = append(artistNames, artist.Name)
	}

	newReleaseYear, newReleaseYearError := wikipedia.GetTrackReleaseYear("Hotel California", artistNames)
	log.Println("NEW RELEASE YEAR", newReleaseYear)
	if newReleaseYearError != nil {
		return newReleaseYearError
	}

	oldestReleaseYear := uint(math.Min(float64(newReleaseYear), float64(track.ReleaseYear)))
	log.Println("OLDEST RELEASE YEAR", oldestReleaseYear)
	if oldestReleaseYear != track.ReleaseYear {
		database.Model(&track).Updates(models.Track{
			ReleaseYear: oldestReleaseYear,
		})

		return nil
	}

	return nil
}
