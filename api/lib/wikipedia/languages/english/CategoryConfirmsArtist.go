package wikipedia_languages_english

import (
	wikipediaModels "api/lib/wikipedia/models"
	"regexp"
)

func CategoryConfirmsArtist(category wikipediaModels.Category, artistName string) bool {
	regex := regexp.MustCompile("(?i)Category:\\s*" + regexp.QuoteMeta(artistName) + "(\\s+\\((band|group)\\))?\\s+songs")
	match := regex.FindStringSubmatch(category.Title)

	return match != nil
}
