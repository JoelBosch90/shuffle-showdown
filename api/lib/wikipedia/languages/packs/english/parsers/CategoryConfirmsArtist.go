package parsers

import (
	"api/lib/wikipedia/models"
	"regexp"
)

func CategoryConfirmsArtist(category models.Category, artistName string) bool {
	regex := regexp.MustCompile("(?i)Category:\\s*" + regexp.QuoteMeta(artistName) + "(\\s+\\((band|group)\\))?\\s+songs")
	match := regex.FindStringSubmatch(category.Title)

	return match != nil
}
