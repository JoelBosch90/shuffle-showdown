package verification

import (
	models "api/database/models"
	"log"

	"github.com/jinzhu/gorm"
)

func VerifyTrack(database *gorm.DB, trackId string) error {
	log.Println("VERIFYING TRACK", trackId)

	track := models.Track{}
	fetchError := database.Model(&track).Preload("Artists").Where("id = ?", trackId).First(&track).Error
	if fetchError != nil {
		return fetchError
	}

	// artistNames := []string{}
	// for _, artist := range track.Artists {
	// 	artistNames = append(artistNames, artist.Name)
	// }

	// artistNamesString := strings.Join(artistNames, ", ")

	return nil
}
