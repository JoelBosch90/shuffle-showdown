package verification

import (
	models "api/database/models"
	"api/lib/wikipedia"
	"log"
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

	newReleaseYear, _ := wikipedia.GetTrackReleaseYear(track.Name, artistNames)
	oldestReleaseYear := uint(math.Min(float64(newReleaseYear), float64(track.ReleaseYear)))
	log.Println("TRACK", track.Name, "ARTISTS", artistNames, "RELEASE YEAR", track.ReleaseYear)
	log.Println("NEW RELEASE YEAR", newReleaseYear)
	log.Println("OLDEST RELEASE YEAR", oldestReleaseYear)
	if oldestReleaseYear != track.ReleaseYear {
		return database.Model(&track).Updates(models.Track{
			ReleaseYear: oldestReleaseYear,
			VerifiedAt:  time.Now(),
		}).Error
	}

	return nil
}
