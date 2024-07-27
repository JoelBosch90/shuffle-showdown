package verification

import (
	models "api/database/models"
	"api/lib/helpers"
	"api/lib/wikipedia"

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

	newReleaseYear, newReleaseMonth, newReleaseDay, newReleaseDateError := wikipedia.GetTrackReleaseDate(track.Name, artistNames)
	if newReleaseDateError != nil {
		return newReleaseDateError
	}

	if helpers.IsOlderDateThan(newReleaseYear, newReleaseMonth, newReleaseDay, track.ReleaseYear, track.ReleaseMonth, track.ReleaseDay) {
		database.Model(&track).Updates(models.Track{
			ReleaseYear:  newReleaseYear,
			ReleaseMonth: newReleaseMonth,
			ReleaseDay:   newReleaseDay,
		})

		return nil
	}

	return nil
}
