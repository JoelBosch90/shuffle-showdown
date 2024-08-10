package parsers

import (
	"api/lib/wikipedia/models"
	"regexp"
)

func CategoryConfirmsArtist(category models.Category, artistName string) bool {
	sanitizedArtistName := regexp.QuoteMeta(artistName)
	regex := regexp.MustCompile("(?i)Catégorie:\\s*(chanson de|Chanson interprétée par) .*" + sanitizedArtistName + ".*")
	match := regex.FindStringSubmatch(category.Title)

	return match != nil
}
